import React, { useEffect, useState } from 'react';
import { useApp } from './AppContext';
import logo from "./assets/logo.svg";

const Auth = () => {
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
                    overflow-y-auto focus:outline-none w-[22rem]">
                    <div className='flex items-center mb-4 gap-4'>
                        <img className='size-12 border-2 bg-figma-white border-figma-black p-1 rounded-lg' src={logo} />
                        <h2 className='text-2xl text-figma-black'>Login</h2>
                    </div>
                    <div className=''>
                        <input className="w-80 border border-figma-gray bg-figma-white focus:outline-none rounded-lg p-2"
                            placeholder="Email"
                            onChange={(e) => { setNewGameUrl(e.target.value) }}
                        />
                    </div>

                    <div className='mt-4'>
                        <input className="w-80 border border-figma-gray bg-figma-white focus:outline-none rounded-lg p-2"
                            type='password'
                            placeholder="Password"
                            onChange={(e) => { setNewGameUrl(e.target.value) }}
                        />
                    </div>
                    {/*<button className='rounded-lg mt-2 text-xs underline ml-1 text-figma-black'>
                        Reset password
                    </button>*/}


                    <div className='flex gap-4 justify-between  mt-5'>
                        <button className='bg-figma-green rounded-lg p-2 px-4 text-figma-white font-bold'>
                            Login
                        </button>
                        <button className='rounded-lg p-2 px-4 text-figma-black'>
                            Register
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Auth;
