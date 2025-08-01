import settings from "./assets/settings.svg"
import Dropdown from "./dropdown";
import { useApp, CurrSheet, ColSuffix } from "./AppContext";

const BottomBar = () => {
    const {
        currSheet, setCurrSheet,
        setCurrTable,
        columns, setColumns,
        setSheetModal,
        setSettingsModal,
        loginData,
    } = useApp();

    return (
        <div className="flex flex-row gap-4 items-center">
            <button className="hover:scale-110 transition-transform duration-100 mr-1"
                onClick={() => setSettingsModal(true)}
            >
                <img className="" src={settings} />
            </button>
            <Dropdown
                options={(loginData?.opened_sheet?.sheets_id_names ?? []).map(item => ({
                    value: item.id,
                    label: item.name
                }))}
                placeholder={loginData?.opened_sheet?.name != ""
                    ? loginData?.opened_sheet?.name
                    : "Sheets"
                }
                onSelect={(item) => {
                    localStorage.setItem(currSheet + ColSuffix, JSON.stringify(columns));
                    setCurrSheet(item.value)
                    localStorage.setItem(CurrSheet, item.value);
                    setColumns(() => {
                        const stored = localStorage.getItem(item.value + ColSuffix);
                        return stored ? JSON.parse(stored) : [];
                    })
                    setCurrTable(() => {
                        const stored = localStorage.getItem(item.value);
                        return stored ? JSON.parse(stored) : [];
                    });
                }}
                isDown={false}
                addNewValue={() => setSheetModal(true)}
            />
        </div>
    );
}

export default BottomBar;

