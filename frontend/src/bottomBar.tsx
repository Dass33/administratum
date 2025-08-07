import settings from "./assets/settings.svg"
import Dropdown, { DropdownOption } from "./dropdown";
import { useApp, Sheet, EnumSheetTypes } from "./AppContext";
import { NewItemProps } from "./NewItemModal.tsx";
import { useEffect, useState } from "react";


const SheetTypesOptions = [{ value: EnumSheetTypes.LIST, label: "Questions" },
{ value: EnumSheetTypes.MAP, label: "Config" }]

const BottomBar = () => {
    const {
        currSheet, setCurrSheet,
        setColumns,
        setNewItemModal,
        setSettingsModal,
        accessToken,
        setSheetDeleted,
        currTable,
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
        setSheetDeleted(false);
    }
    const selectSheets = (item: DropdownOption) => {
        getCurrSheet(item.value, accessToken ?? "", setData)
    }

    const [sheetType, setSheetType] = useState(EnumSheetTypes.LIST);

    const addNewValue = (setSelected: Function) => {
        const props: NewItemProps = {
            currNames: currSheet?.sheets_id_names ?? [],
            assignNewName: (name: string) => {
                if (!currSheet) return
                createSheet(name, sheetType, currSheet.branch_id_name.id, accessToken, (sheet: Sheet) => {
                    setData(sheet);
                    setSelected({ value: sheet.id, label: sheet.name });
                });
            },
            ExpandingComponent: () => (
                <SetSheetType setData={setSheetType} />
            )
        }
        setNewItemModal(props)
        setSheetType(EnumSheetTypes.LIST)
    }

    const assignNewName = (name: string, option: DropdownOption, setSelected: Function) => {
        if (!currSheet) return
        renameSheet(name, currSheet.branch_id_name.id, accessToken);

        const newSheetNames = currSheet.sheets_id_names.map(idName => {
            if (idName.id === option.value) {
                setSelected({ value: idName.id, label: name });
                return { id: idName.id, name: name };
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

    const updateValue = (option: DropdownOption, setSelected: Function) => {
        const idName = { name: option.label, id: option.value }
        const props: NewItemProps = {
            currNames: currSheet?.sheets_id_names ?? [],
            defaultIdName: idName,
            assignNewName(name: string) { assignNewName(name, option, setSelected) },
            deleteItem() { delteItem(option) },
        }
        setNewItemModal(props)
    }

    const everyRender = (setSelected: Function) => {
        useEffect(() => {
            if (!currSheet) return
            setSelected({ name: currSheet.name, value: currSheet.id })
        }, [currTable])
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
                updateValue={updateValue}
                everyRender={everyRender}
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

const createSheet = (name: string, sheetType: string, branchId: string, token: string | undefined, setData: Function) => {
    const createSheetParams: { Name: string, Type: string, BranchID: string } = {
        Name: name,
        BranchID: branchId,
        Type: sheetType,
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
            console.log("recieved", result.type)
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


const SetSheetType: React.FC<{ setData: Function }> = ({ setData }) => (
    <div className='flex justify-between items-ceter my-4'>
        <h2 className="text-xl mr-4 font-medium my-auto text-figma-black">Sheet Type</h2>
        <Dropdown
            options={SheetTypesOptions}
            placeholder={"Select Type"}
            defaultValue={SheetTypesOptions[0].value}
            onSelect={(e) => setData(e.value)}
        />
    </div>
);

export default BottomBar;
