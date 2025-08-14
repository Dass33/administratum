import React, { useEffect, useState } from 'react';
import { Domain, NullString, TableData, useApp } from './AppContext';
import danger from "./assets/danger.svg";

const SettingsModal = () => {
    const {
        setSettingsModal,
        gameUrl, setGameUrl,
        setAuthenticated,
        setColumns,
        setSheets,
        setProjectName,
        setCurrBranch,
        setCurrSheet,
        accessToken,
        currTable
    } = useApp();
    const stopPropagation = (e: React.MouseEvent) => {
        e.stopPropagation();
    };

    const [newGameUrl, setNewGameUrl] = useState(gameUrl);
    const [logutError, SetLogoutError] = useState(false);

    const saveAndExit = () => {
        setGameUrl(newGameUrl);
        changeGameUrl(newGameUrl, currTable, accessToken ?? "");
        setSettingsModal(false);
    }

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'Escape') saveAndExit()
        };

        document.addEventListener('keydown', handleKeyDown);
        return () => {
            document.removeEventListener('keydown', handleKeyDown);
        };
    }, [setSettingsModal]);

    const handleLogout = async () => {
        await fetch(Domain + "/logout", {
            method: "POST",
            credentials: "include"
        })
            .then(response => {
                if (response.status != 204) throw "Logout error"
            })
            .then(() => {
                setAuthenticated(false);
                setSettingsModal(false);
                setColumns([]);
                setSheets([]);
                setProjectName(undefined);
                setCurrBranch(undefined);
                setCurrSheet(undefined);
            })
            .catch(err => {
                SetLogoutError(true);
                console.error(err);
            })
    };

    return (
        <div
            className="fixed inset-0 bg-black bg-opacity-15 flex justify-center items-center z-50"
            onClick={saveAndExit}
        >
            <div onClick={stopPropagation}>
                <div className="p-4 rounded-lg text-figma-black bg-figma-white font-medium
                    overflow-y-auto focus:outline-none w-[35rem]">
                    <h2 className='text-2xl mb-4'> Settings </h2>
                    <div className='flex justify-between items-center'>
                        <span>Game url</span>
                        <input className="w-60 border border-figma-gray bg-figma-white focus:outline-none rounded-lg p-2"
                            defaultValue={gameUrl.String}
                            onChange={(e) => {
                                const val = e.target.value
                                setNewGameUrl({ Valid: isValidUrl(val), String: val })
                            }}
                        />
                    </div>

                    <div className='w-full flex justify-end'>
                        <button className='bg-red-600 w-24 rounded-lg p-2 px-4 text-figma-white font-bold mt-4'
                            onClick={() => handleLogout()}
                        >
                            {logutError
                                ? <img className="size-6 mx-auto" src={danger} />
                                : <span>Log Out</span>
                            }
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

const changeGameUrl = (gameUrl: NullString, table: TableData | undefined, token: string) => {
    if (!table) return;
    const postGameUrlParams: { game_url: NullString, table_id: string } = {
        game_url: gameUrl,
        table_id: table.id,
    }

    fetch(Domain + "/change_game_url", {
        method: "PUT",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(postGameUrlParams)
    })
        .then(response => {
            if (response.status < 200 || response.status > 299) {
                throw new Error("Could not change game url");
            }
        })
        .catch(err => {
            console.error(err);
        });
}

function isValidUrl(url: string) {
    try {
        new URL(url);
        return true;
    } catch {
        return false;
    }
}

export default SettingsModal;
