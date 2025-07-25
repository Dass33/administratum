import React from 'react';
import { useApp } from './AppContext';

const CellModal = () => {
    const { setShowCellModal } = useApp();

    const stopPropagation = (e: React.MouseEvent) => {
        e.stopPropagation();
    };

    return (
        <div className="fixed inset-0 bg-black bg-opacity-20 flex justify-center items-center z-50"
            onClick={() => setShowCellModal(false)}>

            <div onClick={stopPropagation}>
                <textarea
                    className="text-figma-black mb-6 bg-figma-white p-6 rounded-lg w-96 h-72 resize-none overflow-y-auto"
                    defaultValue={"Placeholder modal text"}
                />
            </div>
        </div>
    );
};

export default CellModal;
