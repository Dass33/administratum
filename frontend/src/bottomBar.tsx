import settings from "./assets/settings.svg"
import Dropdown, { DropdownOption } from "./dropdown";
import { useApp, Sheet } from "./AppContext";

const BottomBar = () => {
    const {
        setCurrSheet,
        setColumns,
        setSheetModal,
        setSettingsModal,
        accessToken,
        openedSheet,
    } = useApp();

    const optionsSheets = (openedSheet.sheets_id_names ?? []).map(item => ({
        value: item.id,
        label: item.name
    }))
    const placeholderSheets = openedSheet.name != ""
        ? openedSheet.name
        : "Sheets"

    const selectSheets = (item: DropdownOption) => {
        getCurrSheet(item.value, accessToken ?? "", (sheet: Sheet) => {
            setCurrSheet(sheet);
            setColumns(sheet.columns);
        })
    }

    return (
        <div className="flex flex-row gap-4 items-center">
            <button className="hover:scale-110 transition-transform duration-100 mr-1"
                onClick={() => setSettingsModal(true)}
            >
                <img className="" src={settings} />
            </button>
            <Dropdown
                options={optionsSheets}
                placeholder={placeholderSheets}
                onSelect={(item) => selectSheets(item)}
                isDown={false}
                addNewValue={() => setSheetModal(true)}
            />
        </div>
    );
}


const getCurrSheet = (sheet_id: string, token: string, setData: Function) => {
    const url = `/get_sheet/${sheet_id}`;

    fetch(url, {
        method: "GET",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include"
    })
        .then(response => {
            if (response.status !== 200) {
                throw new Error("Could not retrieve sheet");
            }
            return response.json();
        })
        .then((result: Sheet) => {
            setData(result);
        })
        .catch(err => {
            console.error(err);
        });
};

export default BottomBar;
