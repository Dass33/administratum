import { useApp } from "./AppContext";
import Dropdown from "./dropdown";
import { NewNameProps } from "./NewNameModal";

function SelectBranch() {
    const {
        setBranchName,
        currSheet,
        setNewNameModal,
        currTable,
    } = useApp();
    const optionsBranches = (currTable?.branches_names ?? []).map(item => ({
        value: item.id,
        label: item.name
    }))
    const placeholderBranch = currSheet?.branch_id_name.name != ""
        ? currSheet?.branch_id_name.name
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
                    currNames: currSheet?.sheets_id_names ?? [],
                    assignNewName: createBranch,
                }
                setNewNameModal(props)
            }}
        />
    );
}

export default SelectBranch;
