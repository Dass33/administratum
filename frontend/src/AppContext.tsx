import React, { useState, createContext, useContext, useEffect } from "react";
import { NewItemProps } from "./NewItemModal.tsx";

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
    newItemModal: NewItemProps | null,
    setNewItemModal: Function,
    settingsModal: boolean,
    setSettingsModal: Function,
    shareModal: boolean,
    setShareModal: Function,
    gameUrl: NullString,
    setGameUrl: Function,
    projectName: string | undefined,
    setProjectName: Function,
    currBranch: Branch | undefined,
    setCurrBranch: Function,
    authenticated: boolean,
    setAuthenticated: Function,
    accessToken: string | undefined,
    setAccessToken: Function,
    loading: boolean,
    setLoading: Function,
    sheetDeleted: boolean,
    setSheetDeleted: Function,
    tableNames: IdName[],
    setTableNames: Function,
}

export enum EnumColTypes {
    TEXT = 'text',
    NUMBER = 'number',
    BOOL = 'bool',
    EDITION = 'edition',
    ARRAY = 'array',
}

export enum EnumSheetTypes {
    MAP = 'map',
    LIST = 'list',
}

export enum PermissionsEnum {
    OWNER = 'owner',
    CONTRIBUTOR = 'contributor',
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
export const Domain = "localhost:8080"
export const DEFAULT_UUID = "00000000-0000-0000-0000-000000000000"

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
    type: NullString
}

export type Column = {
    name: string
    id: string
    type: string
    required: boolean
    data: ColumnData[]
}

export type Branch = {
    name: string
    id: string
    is_protected: boolean
}

export type Sheet = {
    name: string
    id: string
    type: string
    row_count: number
    columns: Column[]
    branch_id_name: IdName
    sheets_id_names: IdName[]
}

export type TableData = {
    name: string
    id: string
    game_url: NullString
    permision: string
    branches_names: IdName[]
}

export const isValidEmail = (email: string | undefined): boolean => {
    if (!email) return false;
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
};


const AppContext = createContext<AppState | undefined>(undefined);

export const AppProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [cellModal, setCellModal] = useState(null);
    const [currSheet, setCurrSheet] = useState<Sheet | undefined>();
    const [currTable, setCurrTable] = useState<TableData | undefined>();
    const [colModal, setColModal] = useState(-1);
    const [columns, setColumns] = useState<Column[]>([]);

    const [addColumn, setAddColumn] = useState(false);
    const [sheets, setSheets] = useState([]);
    const [newItemModal, setNewItemModal] = useState(null);
    const [settingsModal, setSettingsModal] = useState(false);
    const [shareModal, setShareModal] = useState(false);
    const [gameUrl, setGameUrl] = useState({ Valid: false, String: "" });
    const [projectName, setProjectName] = useState();
    const [currBranch, setCurrBranch] = useState();
    const [authenticated, setAuthenticated] = useState(false);
    const [accessToken, setAccessToken] = useState<string | undefined>();
    const [loading, setLoading] = useState(true);
    const [sheetDeleted, setSheetDeleted] = useState<boolean>(false);
    const [tableNames, setTableNames] = useState<IdName[]>([]);

    useEffect(() => {
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
            .then((data: LoginData) => {
                if (data) {
                    setAuthenticated(true);
                    setLoading(false);
                    setAccessToken(data.token);
                    setCurrSheet(data.opened_sheet)
                    setCurrTable(data.opened_table)
                    setColumns(data.opened_sheet.columns);
                    setTableNames(data.table_names)
                    setGameUrl(data.opened_table.game_url)
                    console.log(data)
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
            newItemModal, setNewItemModal,
            settingsModal, setSettingsModal,
            shareModal, setShareModal,
            gameUrl, setGameUrl,
            projectName, setProjectName,
            currBranch, setCurrBranch,
            authenticated, setAuthenticated,
            accessToken, setAccessToken,
            loading, setLoading,
            sheetDeleted, setSheetDeleted,
            tableNames, setTableNames,
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
