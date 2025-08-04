import settings from "./assets/settings.svg"
import Dropdown, { DropdownOption } from "./dropdown";
import { useApp, Sheet } from "./AppContext";
import { NewNameProps } from "./NewNameModal";

const BottomBar = () => {
    const {
        currSheet, setCurrSheet,
        setColumns,
        setNewNameModal,
        setSettingsModal,
        accessToken,
        setSheetDeleted,
    } = useApp();

    const optionsSheets = (currSheet?.sheets_id_names ?? []).map(item => ({
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
        setSheetDeleted(false);
    }

    const addNewValue = () => {
        const props: NewNameProps = {
            currNames: currSheet?.sheets_id_names ?? [],
            assignNewName: (name: string) => {
                if (!currSheet) return
                createSheet(name, currSheet.branch_id_name.id, accessToken, setData);
            },
        }
        setNewNameModal(props)
    }

    const assignNewName = (name: string, option: DropdownOption) => {
        if (!currSheet) return
        renameSheet(name, currSheet.branch_id_name.id, accessToken);

        const newSheetNames = currSheet.sheets_id_names.map(idName => {
            if (idName.id === option.value) {
                return { id: idName.id, name: name }
            }
            return idName;
        })
        setCurrSheet({
            ...currSheet,
            sheets_id_names: newSheetNames,
        });
    }

    const delteItem = (option: DropdownOption) => {
        deleteSheet(option.value, accessToken)
        if (!currSheet) return;

        const newSheetNames = currSheet.sheets_id_names.filter(
            idName => idName.id !== option.value
        )
        setCurrSheet({
            ...currSheet,
            sheets_id_names: newSheetNames,
        });

        if (currSheet.id === option.value) {
            setSheetDeleted(true);
            return;
        }
    }

    const updateValue = (option: DropdownOption) => {
        const props: NewNameProps = {
            currNames: currSheet?.sheets_id_names ?? [],
            defaultIdName: { name: option.label, id: option.value },
            assignNewName(name: string) { assignNewName(name, option) },
            deleteItem() { delteItem(option) },
        }
        setNewNameModal(props)
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
                addNewValue={addNewValue}
                updateValue={(option) => { updateValue(option) }}
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
            if (response.status < 200 || response.status > 299) {
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

const createSheet = (name: string, branchId: string, token: string | undefined, setData: Function) => {
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
            if (response.status < 200 || response.status > 299) {
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

const renameSheet = (name: string, sheetId: string, token: string | undefined) => {
    const renameSheetParams: { Name: string, SheetId: string } = {
        Name: name,
        SheetId: sheetId,
    }

    fetch("/rename_sheet", {
        method: "PUT",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(renameSheetParams)
    })
        .then(response => {
            if (response.status < 200 || response.status > 299) {
                throw new Error("Could not rename sheet");
            }
        })
        .catch(err => {
            console.error(err);
        });
}

const deleteSheet = (sheetId: string, token: string | undefined) => {
    const deleteSheetParams: { SheetId: string } = {
        SheetId: sheetId,
    };

    fetch('/delete_sheet', {
        method: "DELETE",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(deleteSheetParams)
    })
        .then(response => {
            if (response.status < 200 || response.status > 299) {
                console.log(response.status)
                throw "Could not delete sheet"
            }
        })
        .catch(err => {
            console.error(err);
        });
}

export default BottomBar;
