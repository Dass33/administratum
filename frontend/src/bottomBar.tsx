import settings from "./assets/settings.svg"
import Dropdown, { DropdownOption } from "./dropdown";
import { useApp, Sheet } from "./AppContext";
import { NewNameProps } from "./NewNameModal";

const BottomBar = () => {
    const {
        setCurrSheet,
        setColumns,
        setNewNameModal,
        setSettingsModal,
        accessToken,
        openedSheet,
        currSheet,
    } = useApp();

    const optionsSheets = (openedSheet?.sheets_id_names ?? []).map(item => ({
        value: item.id,
        label: item.name
    }))
    const placeholderSheets = currSheet?.name != ""
        ? currSheet?.name
        : "Sheets"
    const setData = (sheet: Sheet) => {
        setCurrSheet(sheet);
        setColumns(sheet.columns);
    }
    const selectSheets = (item: DropdownOption) => {
        getCurrSheet(item.value, accessToken ?? "", setData)
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
                addNewValue={() => {
                    const props: NewNameProps = {
                        currNames: openedSheet?.sheets_id_names ?? [],
                        assignNewName: (name: string) => {
                            if (!openedSheet) return
                            createSheet(name, openedSheet.branch_id_name.id, accessToken ?? "", setData);
                        },
                    }
                    setNewNameModal(props)
                }}
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

const createSheet = (name: string, branchId: string, token: string, setData: Function) => {
    const createSheetParams: { Name: string, BranchID: string } = {
        Name: name,
        BranchID: branchId,
    }

    fetch("/create_sheet", {
        method: "POST",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(createSheetParams)
    })
        .then(response => {
            if (response.status !== 201) {
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
}

export default BottomBar;
