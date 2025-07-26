import React, { useState, createContext, useContext } from "react";

export type TableType = Record<string, any>[]

interface AppState {
    cellModal: [number, string] | null
    setCellModal: Function
    currTable: TableType
    setCurrTable: Function
    colModal: boolean
    setColModal: Function
    columns: ColumnProps[],
    setColumns: Function,
}

export interface ColumnProps {
    name: string;
    columnType: string;
    required: boolean;
}

const AppContext = createContext<AppState | undefined>(undefined);

export const AppProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const data = [
        { name: "John", age: 30, city: "New York", active: true },
        { name: "Jane", age: 25, city: "Los Angeles", department: "Engineering" },
        { name: "Bob", age: 35, active: false, salary: 75000 },
        { name: "John", age: 30, city: "New York", active: true },
        { name: "Jane", age: 25, city: "Los Angeles", department: "Engineering" },
        { name: "Bob", age: 35, active: false, salary: 75000 },
        { name: "John", age: 30, city: "New York", active: true },
        { name: "Jane", age: 25, city: "Los Angeles", department: "Engineering" },
        { name: "Bob", age: 35, active: false, salary: 75000 },
        { name: "John", age: 30, city: "New York", active: true },
        { name: "Jane", age: 25, city: "Los Angeles", department: "Engineering" },
        { name: "Bob", age: 35, active: false, salary: 75000, questions: "this is a quesetions?", config: "kaldjfl" },
    ];

    const [cellModal, setCellModal] = useState(null);
    const [currTable, setCurrTable] = useState<TableType>(() => {
        const stored = localStorage.getItem('currTable');
        return stored ? JSON.parse(stored) : data;
    });
    const [colModal, setColModal] = useState(false);
    const [columns, setColumns] = useState<ColumnProps[]>([]);

    return (
        <AppContext.Provider value={{
            cellModal, setCellModal,
            currTable, setCurrTable,
            colModal, setColModal,
            columns, setColumns
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
