import { useEffect, useState } from "react";
import { Branch, Domain, useApp, Sheet } from "./AppContext";
import Dropdown, { DropdownOption } from "./dropdown";
import { NewItemProps } from "./NewItemModal.tsx";

type BranchData = {
    Branch: Branch,
    Sheet: Sheet,
}

function SelectBranch() {
    const {
        currBranch, setCurrBranch,
        currSheet,
        setNewItemModal,
        currTable, setCurrTable,
        accessToken,
        setSheetDeleted,
        setCurrSheet,
        setLoading,
        setColumns,
    } = useApp();

    const optionsBranches = (currTable?.branches_id_names ?? []).map(item => ({
        value: item.id,
        label: item.name
    }))
    const placeholderBranch = currSheet?.curr_branch.name != ""
        ? currSheet?.curr_branch.name
        : "Branch"

    const [isProtected, setIsProtected] = useState(false)
    const setData = (res: BranchData) => {
        setCurrBranch(res.Branch);
        setCurrSheet(res.Sheet);
        setLoading(false);
        setColumns(res.Sheet.columns);
        setSheetDeleted(false);
    }

    const addNewValue = () => {
        const props: NewItemProps = {
            currNames: currTable?.branches_id_names ?? [],
            assignNewName: createBranch,
            ExpandingComponent: () => (
                <SetBranchProtection setData={(option: DropdownOption) => {
                    if (option.value == "true") setIsProtected(true);
                    else setIsProtected(false);
                }} />
            )
        }
        setNewItemModal(props)
    }

    const createBranch = (name: string) => {
        postBranch(name, isProtected, currTable?.id, accessToken ?? "", (res: BranchData) => {
            setData(res);
            if (!currTable) return;
            const newBranch = { name: res.Branch.name, id: res.Branch.id }
            setCurrTable({ ...currTable, branches_id_names: [...currTable.branches_id_names, newBranch] })
        })
    }

    const assignNewName = (name: string, option: DropdownOption, setSelected: Function) => {
        // TODO: fix isProtected
        adjustBranch(name, isProtected, option.value, accessToken);

        const newBranchIdNames = currTable?.branches_id_names.map(idName => {
            if (idName.id === option.value) {
                setSelected({ value: idName.id, label: name })
                return { id: idName.id, name: name }
            }
            return idName;
        })

        setCurrBranch({ ...currBranch, name: name, is_protected: isProtected, });
        setCurrTable({ ...currTable, branches_id_names: newBranchIdNames })
    }

    const deleteItem = (option: DropdownOption) => {
        deleteBranch(option.value, accessToken)

        const newBranchIdNames = currTable?.branches_id_names.filter(
            idName => idName.id !== option.value
        )
        setCurrTable({ ...currTable, branches_id_names: newBranchIdNames })

        if (currBranch?.id === option.value) {
            setCurrBranch();
            setSheetDeleted(true);
        }
    }

    const updateValue = (option: DropdownOption, setSelected: Function) => {
        const props: NewItemProps = {
            currNames: currTable?.branches_id_names ?? [],
            defaultIdName: { name: option.label, id: option.value },
            assignNewName(name: string) { assignNewName(name, option, setSelected) },
            deleteItem() { deleteItem(option) },
        }
        setNewItemModal(props)

    }
    const everyRender = (setSelected: Function, currentSelected?: DropdownOption) => {
        useEffect(() => {
            if (!currBranch) return

            const targetSelection = { name: currBranch.name, value: currBranch.id }

            // Only update if the current selection is different
            if (currentSelected?.value !== targetSelection.value) {
                setSelected(targetSelection)
            }
        }, [currBranch?.id, currentSelected?.value])
    }

    return (
        <Dropdown
            options={optionsBranches}
            placeholder={placeholderBranch}
            onSelect={(option) => getCurrBranch(option.value, accessToken ?? "", setData)}
            addNewValue={addNewValue}
            updateValue={updateValue}
            everyRender={everyRender}
        />
    );
}

const SetBranchProtection: React.FC<{ setData: Function }> = ({ setData }) => (
    <div className='flex justify-between items-ceter my-4'>
        <h2 className="text-xl mr-4 font-medium my-auto text-figma-black">Is protected</h2>
        <Dropdown
            options={[{ value: "false", label: "False" },
            { value: "true", label: "True" }]}
            defaultValue={"false"}
            onSelect={(e) => setData(e)}
        />
    </div >
);

export const getCurrBranch = (branch_id: string, token: string, setData: Function) => {
    const url = Domain + `/get_branch/${branch_id}`;

    fetch(url, {
        method: "GET",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include"
    })
        .then(response => {
            if (response.status < 200 || response.status > 299) {
                throw new Error("Could not retrieve branch");
            }
            return response.json();
        })
        .then(result => {
            setData(result);
        })
        .catch(err => {
            console.error(err);
        });
};

const postBranch = (
    name: string,
    isProtected: boolean,
    tableId: string | undefined,
    token: string,
    setData: Function
) => {
    if (!tableId) return;
    const url = Domain + `/create_branch`;
    const postParam: {
        name: string,
        is_protected: boolean,
        table_id: string
    } = {
        name: name,
        is_protected: isProtected,
        table_id: tableId
    }
    fetch(url, {
        method: "POST",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(postParam)
    })
        .then(response => {
            if (response.status < 200 || response.status > 299) {
                throw new Error("Could not create branch");
            }
            return response.json();
        })
        .then(result => {
            setData(result);
        })
        .catch(err => {
            console.error(err);
        });
};

const adjustBranch = (name: string, isProtected: boolean, branchId: string, token: string | undefined) => {
    const adjustParams: { name: string, branch_id: string, is_protected: boolean } = {
        name: name,
        branch_id: branchId,
        is_protected: isProtected,
    }

    fetch(Domain + "/update_branch", {
        method: "PUT",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(adjustParams)
    })
        .then(response => {
            if (response.status < 200 || response.status > 299) {
                throw new Error("Could not rename branch");
            }
        })
        .catch(err => {
            console.error(err);
        });
}

const deleteBranch = (branchId: string, token: string | undefined) => {
    const deleteBranchParams: { branch_id: string } = {
        branch_id: branchId,
    };

    fetch(Domain + '/delete_branch', {
        method: "DELETE",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(deleteBranchParams)
    })
        .then(response => {
            if (response.status < 200 || response.status > 299) {
                console.log(response.status)
                throw "Could not delete branch"
            }
        })
        .catch(err => {
            console.error(err);
        });
}

export default SelectBranch;
