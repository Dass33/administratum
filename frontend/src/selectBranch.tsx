import { useApp } from "./AppContext";
import Dropdown from "./dropdown";
import { NewNameProps } from "./NewNameModal";

function SelectBranch() {
    const {
        setBranchName,
        loginData,
        openedSheet,
        setNewNameModal,
    } = useApp();
    const optionsBranches = (loginData?.opened_table?.branches_names ?? []).map(item => ({
        value: item.id,
        label: item.name
    }))
    const placeholderBranch = loginData?.opened_sheet.branch_id_name.name != ""
        ? loginData?.opened_sheet.branch_id_name.name
        : "Branch"

    const createBranch = (name: string) => {
        console.log(name);
    }

    return (
        <Dropdown
            options={optionsBranches}
            placeholder={placeholderBranch}

            onSelect={(e) => { setBranchName(e.value) }}
            addNewValue={() => {
                const props: NewNameProps = {
                    currNames: openedSheet?.sheets_id_names ?? [],
                    assignNewName: createBranch,
                }
                setNewNameModal(props)
            }}
        />
    );
}

export default SelectBranch;
