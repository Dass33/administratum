import { useApp } from "./AppContext";

export interface JsonTableProps<T = Record<string, any>> {
    data: T[];
}

const Table = <T extends Record<string, any>>({ data }: JsonTableProps<T>) => {
    const { setShowCellModal } = useApp()

    const getAllColumns = (dataArray: T[]): string[] => {
        const columns = new Set<string>();
        dataArray.forEach(obj => {
            if (typeof obj === 'object' && obj !== null) {
                Object.keys(obj).forEach(key => columns.add(key));
            }
        });
        return Array.from(columns);
    };

    const renderCellValue = (value: any): JSX.Element => {
        if (value === null || value === undefined) {
            return <span className="text-gray-400">null</span>;
        }
        if (typeof value === 'object') {
            return <span className="text-blue-600">{JSON.stringify(value)}</span>;
        }
        if (typeof value === 'boolean') {
            return <span className={value ? 'text-green-600' : 'text-red-600'}>{String(value)}</span>;
        }
        if (typeof value === 'number') {
            return <span className="text-purple-600">{value}</span>;
        }
        return <span className="text-gray-800">{String(value)}</span>;
    };

    if (!Array.isArray(data)) {
        return (
            <div className="max-w-full mx-auto">
                <div className="rounded-lg p-6">
                    <div className="p-3 bg-red-100 border border-red-300 text-red-700 rounded-md">
                        Error: Data must be an array of objects
                    </div>
                </div>
            </div>
        );
    }

    if (data.length === 0) {
        return (
            <div className="mx-auto">
                <div className="rounded-lg p-6">
                    <div className="p-3 w-52 bg-yellow-100 border border-yellow-300 text-yellow-700 rounded-md">
                        No data to display
                    </div>
                </div>
            </div>
        );
    }

    const columns = getAllColumns(data);

    return (
        <div className="max-w-full mx-auto">
            <div className="overflow-x-scroll max-h-[calc(100vh-200px)] -mx-5 -my-2 max-w-[calc(100vw-500px)]">
                <table className="table-fixed border-separate border-spacing-5 w-full">
                    <colgroup>
                        {columns.map((_, idx) => (
                            <col key={idx} className="w-40" />
                        ))}
                    </colgroup>
                    <thead>
                        <tr>
                            {columns.map((col, idx) => (
                                <th
                                    key={idx}
                                    className="border border-gray-300 px-2 py-1 text-left font-semibold text-gray-700 truncate"
                                >
                                    {col}
                                </th>
                            ))}
                        </tr>
                        <tr></tr>
                    </thead>
                    <tbody>
                        {data.map((row, rowIdx) => (
                            <tr key={rowIdx}>
                                {columns.map((col, colIdx) => (
                                    <td
                                        key={colIdx}
                                        className="border border-gray-300 px-2 py-1 text-sm truncate"
                                    >
                                        <button
                                            type="button"
                                            className="w-full text-left focus:outline-none"
                                            onClick={() => setShowCellModal(true)}
                                        >
                                            {renderCellValue(row[col])}
                                        </button>
                                    </td>
                                ))}
                            </tr>
                        ))}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default Table;
