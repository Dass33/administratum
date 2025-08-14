import { useEffect, useState } from "react";
import { useApp, ColTypes, EnumColTypes, Column, Sheet, ColumnData, EnumSheetTypes, Domain, PermissionsEnum, Branch } from "./AppContext";
import plus from "./assets/plus.svg";
import cross from "./assets/cross.svg";
import leftArrow from "./assets/left_arrow.svg";

const Table = () => {
    const {
        setCellModal,
        currSheet, setCurrSheet,
        setColModal,
        columns, setColumns,
        setAddColumn,
        sheetDeleted,
        accessToken,
        tableNames,
        currTable,
        currBranch
    } = useApp()

    const [borderColors, setBorderColors] = useState<string[]>([]);
    const [hoveredRow, setHoveredRow] = useState<number | null>(null);
    const [hideTimeout, setHideTimeout] = useState<number | null>(null);
    const [showTimeout, setShowTimeout] = useState<number | null>(null);
    const [deleteButtonPosition, setDeleteButtonPosition] = useState<{ top: number } | null>(null);

    const isConfig = currSheet?.type == EnumSheetTypes.MAP
    const hasPerms = !currBranch?.is_protected || currTable?.permision === PermissionsEnum.OWNER

    const swapColumns = async (currentIndex: number) => {
        if (!columns || currentIndex === 0) return;

        const previousIndex = currentIndex - 1;

        const newColumns = [...columns];
        [newColumns[previousIndex], newColumns[currentIndex]] = [newColumns[currentIndex], newColumns[previousIndex]];
        setColumns(newColumns);
        putSwapColumns(columns[currentIndex].id, columns[previousIndex].id, accessToken)
    };

    useEffect(() => {
        if (!columns) return
        setBorderColors(columns.map(col => {
            const colType = ColTypes.find(item => col.type === item.val)
            if (colType) return colType.color
            return "border-figma-winter"
        }))
    }, [columns]);


    const getMessage = () => {
        if (sheetDeleted) return "Sheet has been deleted.";
        if (tableNames.length == 0) return "Create a new project.";
        if (!columns || !currSheet) return "Loading...";
        return null;
    };

    const message = getMessage();

    if (message || !currSheet) {
        return (
            <div className="h-full flex justify-center items-center">
                <div className="mb-20 font-medium text-2xl text-gray-500/70 rounded-md">
                    {message}
                </div>
            </div >
        );
    }

    const handleDeleteRow = (rowIdx: number) => {
        if (rowIdx == currSheet.row_count) return;
        deleteRow(currSheet, rowIdx, accessToken)

        const newColumns = columns.map(col => {

            const newCol = col
            newCol.data = col.data
                .map(cell => {
                    if (cell.idx <= rowIdx) return cell
                    const newCell = cell;
                    newCell.idx--;
                    return newCell
                })
                .filter(cell => cell.idx != rowIdx);

            return newCol;
        })
        const newSheet = currSheet;
        newSheet.row_count--;
        setCurrSheet(newSheet);
        setColumns(newColumns)

    };

    const handleRowMouseEnter = (rowIdx: number, event: React.MouseEvent<HTMLTableRowElement>) => {
        if (rowIdx == currSheet.row_count) return
        if (hideTimeout) {
            clearTimeout(hideTimeout);
            setHideTimeout(null);
        }
        if (showTimeout) {
            clearTimeout(showTimeout);
            setShowTimeout(null);
        }

        const row = event.currentTarget;
        const scrollContainer = row.closest('.overflow-x-scroll');
        if (scrollContainer && row) {
            const rowRect = row.getBoundingClientRect();
            const containerRect = scrollContainer.getBoundingClientRect();
            const offset = 9;
            const relativeTop = rowRect.top - containerRect.top - offset;

            const timeout = setTimeout(() => {
                setDeleteButtonPosition({ top: relativeTop });
                setHoveredRow(rowIdx);
                setShowTimeout(null);
            }, 150);

            setShowTimeout(timeout);
        }
    };

    const handleRowMouseLeave = () => {
        const timeout = setTimeout(() => {
            setHoveredRow(null);
            setDeleteButtonPosition(null);
        }, 600);
        setHideTimeout(timeout);
    };

    const handleDeleteButtonMouseEnter = () => {
        if (hideTimeout) {
            clearTimeout(hideTimeout);
            setHideTimeout(null);
        }
    };

    const handleDeleteButtonMouseLeave = () => {
        const timeout = setTimeout(() => {
            setHoveredRow(null);
            setDeleteButtonPosition(null);
        }, 300);
        setHideTimeout(timeout);
    };

    return (
        <div className="ml-2 max-w-full flex items-start">
            <div className="overflow-x-scroll max-h-[calc(100vh-200px)] -mx-5 -my-2 md:max-w-[65vw] xl:max-w-[calc(100vw-500px)]">
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
                                    <div className="flex justify-between">
                                        <button
                                            type="button"
                                            disabled={isConfig || !hasPerms}
                                            className="w-full text-left focus:outline-none"
                                            onClick={() => {
                                                if (!isConfig && hasPerms) setColModal(idx)
                                            }}
                                        >
                                            {col.name}
                                        </button>

                                        <button
                                            className={`${!isConfig && hasPerms && idx ? "block" : "hidden"}`}
                                            onClick={() => swapColumns(idx)}>
                                            <img className="size-7" src={leftArrow} alt="left arrow" />
                                        </button>
                                    </div>
                                </th>
                            ))}
                        </tr>
                        <tr></tr>
                    </thead>
                    <tbody>
                        {Array.from({ length: (currSheet?.row_count ?? 0) + Number(hasPerms) }, (_, rowIdx) => (
                            <tr key={rowIdx}
                                onMouseEnter={(e) => handleRowMouseEnter(rowIdx, e)}
                                onMouseLeave={handleRowMouseLeave}
                            >
                                {columns.map((col, colIdx) => (
                                    <td key={colIdx}
                                        className={`border px-3 py-2 text-sm truncate
                                        ${rowIdx == currSheet.row_count || validateCellType(rowIdx, col, currBranch)
                                                ? "border-gray-300"
                                                : "border-red-600"}`}
                                    >
                                        <button
                                            type="button"
                                            className="w-full text-left focus:outline-none"
                                            disabled={!hasPerms}
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
            <div className={`ml-4 flex flex-col items-center justify-start relative
                            ${!hasPerms && "hidden"}`}>

                <button className="w-12 h-12 flex items-center justify-center text-[3rem] font-light hover:scale-125 transition-transform duration-100 flex-shrink-0"
                    disabled={isConfig}
                    onClick={() => {
                        setColModal(columns.length)
                        setAddColumn(true)
                    }}>
                    <img src={plus} className={`${isConfig && "hidden"} w-7 h-7`} />
                </button>

                {hoveredRow !== null && deleteButtonPosition && (
                    <div
                        className="absolute w-10 h-10 flex items-center justify-center hover:scale-125 transition-all duration-100"
                        style={{
                            top: `${deleteButtonPosition.top}px`,
                        }}
                        onMouseEnter={handleDeleteButtonMouseEnter}
                        onMouseLeave={handleDeleteButtonMouseLeave}
                    >
                        <button
                            className="w-full h-full flex items-center justify-center"
                            onClick={() => handleDeleteRow(hoveredRow)}
                        >
                            <img src={cross} className="w-7 h-7" />
                        </button>
                    </div>
                )}
            </div>
        </div>
    );
};

const renderCellValue = (data: ColumnData[], idx: number): JSX.Element => {
    const item = data.find(item => item.idx == idx)
    if (!item) {
        return <span className="text-gray-400">null</span>;
    }
    const value = item.value.String
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

const validateCellType = (idx: number, col: Column, currBranch: Branch | undefined): boolean => {
    const item = col.data.find(item => item.idx == idx)
    if (!item || !item.value.Valid) {
        return !col.required;
    }

    const val = item.value.String

    // Check if it's an enum type
    const baseTypes = [EnumColTypes.TEXT, EnumColTypes.NUMBER, EnumColTypes.BOOL, EnumColTypes.ARRAY];
    const isEnum = !baseTypes.includes(col.type as EnumColTypes);

    if (isEnum) {
        const enumItem = currBranch?.enums?.find((e) => e.name === col.type);
        if (!enumItem) {
            return true;
        }
        return enumItem.vals.includes(val);
    }

    switch (col.type) {
        case EnumColTypes.TEXT:
            return typeof val === 'string';

        case EnumColTypes.NUMBER: {
            const numVal = Number(val);
            return !isNaN(numVal) && isFinite(numVal);
        }

        case EnumColTypes.BOOL: {
            const lowerVal = val.toLowerCase().trim();
            return lowerVal === 'true' || lowerVal === 'false';
        }

        case EnumColTypes.ARRAY:
            try {
                const parsedArray = JSON.parse(val);
                if (!Array.isArray(parsedArray)) return false;

                if (col.required) {
                    return parsedArray.length > 0 && parsedArray.some((item) =>
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

const deleteRow = (sheet: Sheet, rowIdx: number, token: string | undefined) => {
    const deleteRowParams: { sheet_id: string, row_idx: number } = {
        sheet_id: sheet.id,
        row_idx: rowIdx,
    };

    fetch(Domain + '/delete_row', {
        method: "DELETE",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(deleteRowParams)
    })
        .then(response => {
            if (response.status < 200 || response.status > 299) {
                throw "Could not delete row"
            }
        })
        .catch(err => {
            console.error(err);
        });
}

const putSwapColumns = (columnId1: string, columnId2: string, token: string | undefined) => {
    const swapColsParams: { column_id1: string, column_id2: string } = {
        column_id1: columnId1,
        column_id2: columnId2,
    };

    fetch(Domain + '/swap_columns', {
        method: "PUT",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(swapColsParams)
    })
        .then(response => {
            if (response.status < 200 || response.status > 299) {
                throw "Could not swap columns"
            }
        })
        .catch(err => {
            console.error(err);
        });
}

export default Table;
