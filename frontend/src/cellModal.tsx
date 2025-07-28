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

    const updateCell = (newVal: string, rowIndex: number, col: string) => {
        setCurrTable((prevTable: TableType) => {
            const newTable = [...prevTable];
            newTable[rowIndex][col] = newVal;
            return newTable;
        });
    }

    const removeEmptyRow = (newVal: string, rowIndex: number, col: string): boolean => {
        if (currTable.length - 1 == rowIndex) return false
        const len = Object.entries(currTable[rowIndex]).filter(([key, val]) => {
            return key != col && val
        }).length
        if (!len && !newVal) {
            setCurrTable(currTable.filter((_, idx) => { return rowIndex != idx }))
            return true
        }
        return false
    }

    const saveAndExit = () => {
        setCellModal(null);
        if (!cellModal) return
        const updatedValue = textareaRef.current?.value ?? cellVal;
        const rowIndex = cellModal[0];
        const col = cellModal[1];
        if (!removeEmptyRow(updatedValue, rowIndex, col)) {
            updateCell(updatedValue, rowIndex, col)
        }
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
