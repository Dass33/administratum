import { useEffect, useState } from "react";
import { useApp, ColTypes, EnumColTypes, Column, ColumnData } from "./AppContext";
import plus from "./assets/plus.svg";

const Table = () => {
    const {
        setCellModal,
        currSheet,
        setColModal,
        columns,
        setAddColumn,
        sheetDeleted,
    } = useApp()
    // console.log(columns)

    const [borderColors, setBorderColors] = useState<string[]>([]);

    useEffect(() => {
        if (!columns) return
        setBorderColors(columns.map(col => {
            const colType = ColTypes.find(item => col.type === item.val)
            if (colType) return colType.color
            return ColTypes[0].color
        }))
    }, [columns]);

    if (!columns || !currSheet) {
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

    if (sheetDeleted) {
        return (
            <div className="max-w-full mx-auto">
                <div className="rounded-lg p-6">
                    <div className="p-3 text-figma-black rounded-md">
                        Sheet has been deleted.
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
                        {Array.from({ length: (currSheet?.row_count ?? 0) + 1 }, (_, rowIdx) => (
                            <tr key={rowIdx}>
                                {columns.map((col, colIdx) => (
                                    <td key={colIdx}
                                        className={`border px-3 py-2 text-sm truncate
                                        ${rowIdx == currSheet.row_count || validateCellType(rowIdx, col)
                                                ? "border-gray-300"
                                                : "border-red-600"}`}
                                    >
                                        <button
                                            type="button"
                                            className="w-full text-left focus:outline-none"
                                            onClick={() => setCellModal([rowIdx, col])}
                                        >
                                            {renderCellValue(col.data, rowIdx)}
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

const renderCellValue = (data: ColumnData[], idx: number): JSX.Element => {
    const item = data.find(item => item.idx == idx)
    if (!item) {
        return <span className="text-gray-400">null</span>;
    }
    let value = item.value.String
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

const validateCellType = (idx: number, col: Column): boolean => {
    const item = col.data.find(item => item.idx == idx)
    if (!item || !item.value.Valid) {
        return !col.required;
    }

    const val = item.value.String

    switch (col.type) {
        case EnumColTypes.TEXT:
            return typeof val === 'string';

        case EnumColTypes.NUMBER:
            const numVal = Number(val);
            return !isNaN(numVal) && isFinite(numVal);

        case EnumColTypes.BOOL:
            const lowerVal = val.toLowerCase().trim();
            return lowerVal === 'true' || lowerVal === 'false';

        case EnumColTypes.ARRAY:
            try {
                const parsedArray = JSON.parse(val);
                if (!Array.isArray(parsedArray)) return false;

                if (col.required) {
                    return parsedArray.length > 0 && parsedArray.some((item: any) =>
                        item !== null && item !== undefined && item !== ""
                    );
                }
                return true;
            } catch {
                return false;
            }

        default:
            return false;
    }
};

export default Table;
