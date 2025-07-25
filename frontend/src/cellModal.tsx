import React, { useEffect, useRef, useState } from 'react';
import { useApp, TableType } from './AppContext';

const CellModal = () => {
    const { setCellModal, cellModal, currTable, setCurrTable } = useApp();
    const textareaRef = useRef<HTMLTextAreaElement>(null);

    const stopPropagation = (e: React.MouseEvent) => {
        e.stopPropagation();
    };
    const initCellVal = cellModal
        ? currTable[cellModal[0]][cellModal[1]]
        : ""
    const [cellVal, setCellVal] = useState(initCellVal)

    const saveAndExit = () => {
        if (cellModal) {
            const updatedValue = textareaRef.current?.value ?? cellVal;
            setCurrTable((prevTable: TableType) => {
                const newTable = [...prevTable];
                const rowIndex = cellModal[0];
                const col = cellModal[1];

                newTable[rowIndex][col] = updatedValue;

                return newTable;
            });
        }
        setCellModal(null);
    }

    useEffect(() => {
        const textarea = textareaRef.current;
        if (textarea) {
            textarea.focus();
            const length = textarea.value.length;
            textarea.setSelectionRange(length, length);
        }
        const handleKeyDown = (e: KeyboardEvent) => {
            if (e.key === 'Escape') saveAndExit()
        };

        document.addEventListener('keydown', handleKeyDown);

        return () => {
            document.removeEventListener('keydown', handleKeyDown);
        };
    }, [setCellModal]);

    return (
        <div
            className="fixed inset-0 bg-black bg-opacity-15 flex justify-center items-center z-50"
            onClick={saveAndExit}
        >
            <div onClick={stopPropagation}>
                <textarea
                    ref={textareaRef}
                    className="text-figma-black mb-6 bg-figma-white p-6 rounded-lg w-[35rem] h-72 resize-none overflow-y-auto focus:outline-none"
                    defaultValue={cellVal}
                    onChange={(e) => setCellVal(e.target.value)}
                />
            </div>
        </div>
    );
};

export default CellModal;
