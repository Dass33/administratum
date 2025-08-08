import React, { useState, useRef, useCallback, useEffect } from 'react';
import BottomBar from "./bottomBar";
import GameView from "./gameView"
import ShareButton from "./shareButton";
import SelectBranch from "./selectBranch";
import SelectProject from "./selectProject";
import Permissions from "./permissions";
import Table from "./table";
import { useApp } from "./AppContext";
import CellModal from "./cellModal";
import ColModal from "./colModal";
import NewItemModal from "./NewItemModal.tsx";
import SettingsModal from "./settingsModal";
import ShareModal from "./shareModal";
import Auth from "./auth";
import RefershButton from './RefreshButton.tsx';

function App() {
    const {
        cellModal,
        colModal,
        newItemModal,
        settingsModal,
        shareModal,
        authenticated,
        loading,
    } = useApp()

    const [leftWidth, setLeftWidth] = useState(70);
    const [isDragging, setIsDragging] = useState(false);
    const containerRef = useRef<HTMLDivElement>(null);

    const handleMouseDown = useCallback((e: React.MouseEvent<HTMLDivElement>) => {
        setIsDragging(true);
        e.preventDefault();
    }, []);

    const handleMouseMove = useCallback((e: MouseEvent) => {
        if (!isDragging || !containerRef.current) return;

        const containerRect = containerRef.current.getBoundingClientRect();
        const newLeftWidth = ((e.clientX - containerRect.left) / containerRect.width) * 100;

        const clampedWidth = Math.min(Math.max(newLeftWidth, 20), 90);
        setLeftWidth(clampedWidth);
    }, [isDragging]);

    const handleMouseUp = useCallback(() => {
        setIsDragging(false);
    }, []);

    useEffect(() => {
        if (isDragging) {
            document.addEventListener('mousemove', handleMouseMove);
            document.addEventListener('mouseup', handleMouseUp);
            document.body.style.cursor = 'col-resize';
            document.body.style.userSelect = 'none';

            return () => {
                document.removeEventListener('mousemove', handleMouseMove);
                document.removeEventListener('mouseup', handleMouseUp);
                document.body.style.cursor = '';
                document.body.style.userSelect = '';
            };
        }
    }, [isDragging, handleMouseMove, handleMouseUp]);

    if (loading) return (
        <div className='bg-figma-white h-screen pl-10 xl:pr-5 flex items-center'>
            <h1 className="text-figma-black text-2xl text-center w-full">Loading...</h1>
        </div>
    )

    if (!authenticated) return (
        <div className='bg-figma-white h-screen pl-10 xl:pr-5 flex flex-row'>
            <Auth />
        </div>
    )

    return (
        <div className='bg-figma-white h-screen px-2 md:px-9 flex flex-row' ref={containerRef}>
            {cellModal && <CellModal />}
            {colModal > -1 && <ColModal />}
            {newItemModal && <NewItemModal {...newItemModal} />}
            {settingsModal && <SettingsModal />}
            {shareModal && <ShareModal />}

            <div className="flex flex-col my-10 overflow-hidden !w-full !md:w-auto"
                style={{ width: `${leftWidth}%` }}
            >
                <div className="flex flex-row gap-4 min-w-[30rem]">
                    <ShareButton />
                    <div className="flex items-center text-lg">
                        <SelectProject />
                        <span className="mx-2">/</span>
                        <SelectBranch />
                    </div>
                    <Permissions />
                    <RefershButton />
                </div>
                <div className="mt-10 grow overflow-hidden">
                    <Table />
                </div>
                <div className="bottom-1 absolute my-2">
                    <BottomBar />
                </div>
            </div>

            <div className="w-4 md:flex items-center justify-center self-center cursor-col-resize group hidden"
                style={{ height: '256px' }}
                onMouseDown={handleMouseDown}
            >
                <div className={`
            w-1.5 h-24 bg-gray-300 rounded-full transition-colors duration-200
            group-hover:bg-gray-400
            ${isDragging ? 'bg-gray-400' : ''}
        `}
                ></div>
            </div>

            <div className="overflow-hidden hidden md:block"
                style={{ width: `${100 - leftWidth}%` }}
            >
                <GameView />
            </div>
        </div>
    )
}

export default App
