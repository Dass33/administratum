import React, { useState, createContext, useContext, useEffect } from "react";

interface AppState {
    cellModal: [number, Column] | null
    setCellModal: Function
    currTable: TableData | undefined
    setCurrTable: Function
    colModal: number
    setColModal: Function
    columns: Column[],
    setColumns: Function,
    addColumn: boolean,
    setAddColumn: Function,
    sheets: string[],
    setSheets: Function,
    currSheet: Sheet | undefined,
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
    loginData: LoginData | undefined,
    setLoginData: Function,
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

export type IdName = {
    name: string
    id: string
}

export type NullString = {
    String: string
    Valid: boolean
}

export type LoginData = {
    email: string
    token: string
    opened_table: TableData
    opened_sheet: Sheet
    table_names: IdName[]
}

export type ColumnData = {
    id: string
    idx: number
    value: NullString
}

export type Column = {
    name: string
    id: string
    type: string
    required: boolean
    data: ColumnData[]
}

export type Sheet = {
    name: string
    id: string
    row_count: number
    columns: Column[]
    branch_id_name: IdName
    sheets_id_names: IdName[]
}

export type TableData = {
    name: string
    id: string
    game_url: string
    permision: string
    branches_names: IdName[]
}

const AppContext = createContext<AppState | undefined>(undefined);

export const AppProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [cellModal, setCellModal] = useState(null);
    const [currSheet, setCurrSheet] = useState<Sheet | undefined>();
    const [currTable, setCurrTable] = useState<TableData | undefined>();
    const [colModal, setColModal] = useState(-1);
    const [columns, setColumns] = useState<Column[]>(() => {
        const stored = localStorage.getItem(currSheet + ColSuffix);
        return stored ? JSON.parse(stored) : [];
    });

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
    const [loginData, setLoginData] = useState<LoginData | undefined>();

    useEffect(() => {
        if (!loginData) return
        if (loginData?.opened_sheet) {
            setCurrSheet(loginData.opened_sheet)
        }
        if (loginData?.opened_sheet?.columns) {
            setColumns(loginData.opened_sheet.columns);
        }
    }, [loginData]);

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
            loginData, setLoginData,
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
