import React, { useState, createContext, useContext } from "react";

export type TableType = Record<string, any>[]

interface AppState {
    cellModal: [number, string] | null
    setCellModal: Function
    currTable: TableType
    setCurrTable: Function
    colModal: number
    setColModal: Function
    columns: ColumnProps[],
    setColumns: Function,
    addColumn: boolean,
    setAddColumn: Function,
}

export interface ColumnProps {
    name: string;
    columnType: string;
    required: boolean;
}

export const ColTypes = [
    { val: 'text', color: "border-figma-stone" },
    { val: 'number', color: "border-figma-rose" },
    { val: 'bool', color: "border-figma-honey" },
    { val: 'edition', color: "border-figma-winter" }
]

const AppContext = createContext<AppState | undefined>(undefined);

export const AppProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    // const data = [
    //     { name: "John", age: 30, city: "New York", active: true },
    //     { name: "Jane", age: 25, city: "Los Angeles" },
    //     { name: "Bob", age: 35, active: false, salary: 75000 },
    // ];
    const data: TableType = [];

    const [cellModal, setCellModal] = useState(null);
    const [currTable, setCurrTable] = useState<TableType>(() => {
        const stored = localStorage.getItem('currTable');
        return stored ? JSON.parse(stored) : data;
    });
    const [colModal, setColModal] = useState(-1);
    const [columns, setColumns] = useState<ColumnProps[]>([]);
    const [addColumn, setAddColumn] = useState(false);

    return (
        <AppContext.Provider value={{
            cellModal, setCellModal,
            currTable, setCurrTable,
            colModal, setColModal,
            columns, setColumns,
            addColumn, setAddColumn,
        }}>
            {children}
        </AppContext.Provider>
    );
};

export const useApp = () => {
    const context = useContext(AppContext);
    if (!context) {
        throw new Error('useApp must be used within a AppProvider');
    }
    return context;
};
