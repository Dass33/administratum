import React, { useState, createContext, useContext, useEffect } from "react";

export type TableType = Record<string, any>[]

interface AppState {
    cellModal: [number, ColumnProps] | null
    setCellModal: Function
    currTable: TableType
    setCurrTable: Function
    colModal: number
    setColModal: Function
    columns: ColumnProps[],
    setColumns: Function,
    addColumn: boolean,
    setAddColumn: Function,
    sheets: string[],
    setSheets: Function,
    currSheet: string,
    setCurrSheet: Function,
    sheetModal: boolean,
    setSheetModal: Function,
    settingsModal: boolean,
    setSettingsModal: Function,
    gameUrl: string,
    setGameUrl: Function,
    projectName: string | undefined,
    setProjectName: Function,
    branchName: string | undefined,
    setBranchName: Function,
    authenticated: boolean,
    setAuthenticated: Function,
    accessToken: string | undefined,
    setAccessToken: Function,
    loading: boolean,
    setLoading: Function,
}

export interface ColumnProps {
    name: string;
    columnType: string;
    required: boolean;
}

export enum EnumColTypes {
    TEXT = 'text',
    NUMBER = 'number',
    BOOL = 'bool',
    EDITION = 'edition',
    ARRAY = 'array',
}

export const ColTypes = [
    { val: EnumColTypes.TEXT, color: "border-figma-stone" },
    { val: EnumColTypes.NUMBER, color: "border-figma-pool" },
    { val: EnumColTypes.BOOL, color: "border-figma-honey" },
    { val: EnumColTypes.EDITION, color: "border-figma-winter" },
    { val: EnumColTypes.ARRAY, color: "border-figma-forest" },
]

export const CurrSheet = 'currSheet'
export const Sheets = 'sheets'
export const ColSuffix = '/columns'

export type Column = {
    name: string
    type: string
    require: boolean
    data: any
}

export type Sheet = {
    name: string
    columns: Column
    opened_branch_name: string
}

export type TableData = {
    game_url: string
    permision: string
    opened_sheet: Sheet
    sheets_names: string[]
    branches_names: string[]
}

const AppContext = createContext<AppState | undefined>(undefined);

export const AppProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    // const data = [
    //     { name: "John", age: 30, city: "New York", active: true },
    //     { name: "Jane", age: 25, city: "Los Angeles" },
    //     { name: "Bob", age: 35, active: false, salary: 75000 },
    // ];
    const data: TableType = [];

    const [cellModal, setCellModal] = useState(null);
    const [currSheet, setCurrSheet] = useState(() => {
        const stored = localStorage.getItem(CurrSheet);
        return stored ? stored : "config";
    });
    const [currTable, setCurrTable] = useState<TableType>(() => {
        const stored = localStorage.getItem(currSheet);
        return stored ? JSON.parse(stored) : data;
    });
    const [colModal, setColModal] = useState(-1);
    const [columns, setColumns] = useState<ColumnProps[]>(() => {
        const stored = localStorage.getItem(currSheet + ColSuffix);
        return stored ? JSON.parse(stored) : [];
    });

    useEffect(() => {
        if (columns.length) return
        fetch('http://localhost:8080/columns')
            .then(response => response.json())
            .then(data => {
                setColumns(data);
                localStorage.setItem(currSheet + ColSuffix, JSON.stringify(data));
            })
            .catch(_ => {
                setColumns([]);
            });
    }, [currSheet]);

    const [addColumn, setAddColumn] = useState(false);
    const [sheets, setSheets] = useState([]);
    const [sheetModal, setSheetModal] = useState(false);
    const [settingsModal, setSettingsModal] = useState(false);
    const [gameUrl, setGameUrl] = useState("https://dass33.github.io/guess_game/");
    const [projectName, setProjectName] = useState();
    const [branchName, setBranchName] = useState();
    const [authenticated, setAuthenticated] = useState(false);
    const [accessToken, setAccessToken] = useState<string | undefined>();
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        localStorage.setItem(currSheet, JSON.stringify(currTable));
    }, [currTable]);

    useEffect(() => {
        interface Token {
            token: string
        }
        fetch('/refresh', {
            method: "POST",
            credentials: "include"
        })
            .then(response => {
                if (response.status != 200) {
                    throw "Not valid refresh token"
                }
                return response.json()
            })
            .then((data: Token) => {
                if (data) {
                    setAccessToken(data.token);
                    setAuthenticated(true);
                    setLoading(false);
                }
            })
            .catch(err => {
                setAuthenticated(false);
                setLoading(false);
                console.error(err);
            });
    }, []);

    return (
        <AppContext.Provider value={{
            cellModal, setCellModal,
            currTable, setCurrTable,
            colModal, setColModal,
            columns, setColumns,
            addColumn, setAddColumn,
            sheets, setSheets,
            currSheet, setCurrSheet,
            sheetModal, setSheetModal,
            settingsModal, setSettingsModal,
            gameUrl, setGameUrl,
            projectName, setProjectName,
            branchName, setBranchName,
            authenticated, setAuthenticated,
            accessToken, setAccessToken,
            loading, setLoading,
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
