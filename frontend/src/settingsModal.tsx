import React, { useEffect, useState } from 'react';
import { useApp } from './AppContext';
import danger from "./assets/danger.svg";

const SettingsModal = () => {
    const {
        setSettingsModal,
        gameUrl, setGameUrl,
        setAuthenticated,
        setColumns,
        setSheets,
        setProjectName,
        setBranchName,
        setCurrSheet,
    } = useApp();
    const stopPropagation = (e: React.MouseEvent) => {
        e.stopPropagation();
    };

    const [newGameUrl, setNewGameUrl] = useState(gameUrl);
    const [logutError, SetLogoutError] = useState(false);

    const saveAndExit = () => {
        setGameUrl(newGameUrl);
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
        await fetch("/logout", {
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
                setBranchName(undefined);
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
                                if (val == "") return;
                                setNewGameUrl({ Valid: true, String: val })
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

export default SettingsModal;
