import React, { useEffect, useState } from 'react';
import { useApp } from './AppContext';

const SettingsModal = () => {
    const {
        setSettingsModal,
        gameUrl, setGameUrl
    } = useApp();
    const stopPropagation = (e: React.MouseEvent) => {
        e.stopPropagation();
    };

    const [newGameUrl, setNewGameUrl] = useState(gameUrl);

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
                            defaultValue={gameUrl}
                            onChange={(e) => { setNewGameUrl(e.target.value) }}
                        />
                    </div>
                </div>
            </div>
        </div>
    );
};

export default SettingsModal;
