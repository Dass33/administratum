import { useEffect, useState } from "react";
import { useApp, ColTypes, EnumColTypes, ColumnProps } from "./AppContext";
import plus from "./assets/plus.svg";


const Table = () => {
    const {
        setCellModal,
        currTable, setCurrTable,
        setColModal,
        columns,
        setAddColumn,
    } = useApp()

    const [borderColors, setBorderColors] = useState<string[]>([]);

    useEffect(() => {
        if (!columns) return
        setBorderColors(columns.map(col => {
            const colType = ColTypes.find(item => col.columnType === item.val)
            if (colType) return colType.color
            return ColTypes[0].color
        }))
    }, [columns]);



    const isRowEmpty = (row: Record<string, any>) => {
        return columns.every(col => {
            const value = row[col.name];
            return value === null || value === undefined || value === '';
        });
    };

    const createEmptyRow = () => {
        const emptyRow: Record<string, any> = {};
        columns.forEach(col => {
            emptyRow[col.name] = '';
        });
        return emptyRow;
    };
    useEffect(() => {
        if (!columns) return
        if (currTable.length == 0) {
            setCurrTable([...currTable, createEmptyRow()]);
            return
        }
        if (currTable && currTable.length > 0 && columns.length > 0) {
            const lastRow = currTable[currTable.length - 1];
            if (!isRowEmpty(lastRow)) {
                setCurrTable([...currTable, createEmptyRow()]);
            }
        }
    }, [currTable, columns]);

    if (!columns) {
        return (
            <div className="max-w-full mx-auto">
                <div className="rounded-lg p-6">
                    <div className="p-3 text-figma-black rounded-md">
                        Loading...
                    </div>
                </div>
            </div>
        );
    }

    if (!Array.isArray(currTable)) {
        return (
            <div className="max-w-full mx-auto">
                <div className="rounded-lg p-6">
                    <div className="p-3 bg-red-100 border border-red-300 text-red-700 rounded-md">
                        Error: Data malformed
                    </div>
                </div>
            </div>
        );
    }

    return (
        <div className="max-w-full mx-auto flex items-start">
            <div className="overflow-x-scroll max-h-[calc(100vh-200px)] -mx-5 -my-2 max-w-[65vw] xl:max-w-[calc(100vw-500px)]">
                <table className="table-fixed border-separate border-spacing-3 w-full">
                    <colgroup>
                        {columns.map((_, idx) => (
                            <col key={idx} className="w-40" />
                        ))}
                    </colgroup>
                    <thead>
                        <tr>
                            {columns.map((col, idx) => (
                                <th key={idx}
                                    className={`border ${borderColors[idx]} px-3 py-2 text-left font-semibold text-gray-700 truncate`}
                                >
                                    <button
                                        type="button"
                                        className="w-full text-left focus:outline-none"
                                        onClick={() => setColModal(idx)}
                                    >
                                        {col.name}
                                    </button>
                                </th>
                            ))}
                        </tr>
                        <tr></tr>
                    </thead>
                    <tbody>
                        {currTable.map((row, rowIdx) => (
                            <tr key={rowIdx}>
                                {columns.map((col, colIdx) => (
                                    <td key={colIdx}
                                        className={`border px-3 py-2 text-sm truncate
                                        ${!col.required || validateCellType(row[col.name], col)
                                                ? "border-gray-300"
                                                : "border-red-600"}`}
                                    >
                                        <button
                                            type="button"
                                            className="w-full text-left focus:outline-none"
                                            onClick={() => setCellModal([rowIdx, col])}
                                        >
                                            {renderCellValue(row[col.name])}
                                        </button>
                                    </td>
                                ))}
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
            <button className="ml-8 mt-2 text-[3rem] font-light hover:scale-125 transition-transform duration-100"
                onClick={() => {
                    setColModal(true)
                    setAddColumn(true)
                }}>
                <img src={plus} />
            </button>
        </div>
    );
};

const renderCellValue = (value: any): JSX.Element => {
    if (value === null || value === undefined || value === "") {
        return <span className="text-gray-400">null</span>;
    }
    if (typeof value === 'object') {
        return <span>{JSON.stringify(value)}</span>;
    }
    if (typeof value === 'boolean') {
        return <span>{value ? "True" : "False"}</span>;
    }
    if (typeof value === 'number') {
        return <span>{value}</span>;
    }
    return <span>{String(value)}</span>;
}

const validateCellType = (cellVal: any, col: ColumnProps): boolean => {
    if (cellVal === null || cellVal === undefined || cellVal === "") {
        return !col.required;
    }

    switch (col.columnType) {
        case EnumColTypes.TEXT:
            return typeof cellVal === 'string';
        case EnumColTypes.NUMBER:
            return typeof cellVal === 'number' && !isNaN(cellVal);
        case EnumColTypes.BOOL:
            return typeof cellVal === 'boolean';
        case EnumColTypes.ARRAY:
            if (!Array.isArray(cellVal)) return false;
            if (col.required) {
                return cellVal.length > 0 && cellVal.some((item: any) =>
                    item !== null && item !== undefined && item !== ""
                );
            }
            return true;
        default:
            return false;
    }
};

export default Table;
