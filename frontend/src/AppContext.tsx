import React, { useState, createContext, useContext } from "react";

interface AppState {
    showCellModal: boolean
    setShowCellModal: Function
}

const AppContext = createContext<AppState | undefined>(undefined);

export const AppProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [showCellModal, setShowCellModal] = useState(true);

    return (
        <AppContext.Provider value={{
            showCellModal, setShowCellModal
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
