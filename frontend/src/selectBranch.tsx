import { useApp } from "./AppContext";
import Dropdown from "./dropdown";

function SelectBranch() {
    const { setBranchName, loginData } = useApp();
    const optionsBranches = (loginData?.opened_table?.branches_names ?? []).map(item => ({
        value: item.id,
        label: item.name
    }))
    const placeholderBranch = loginData?.opened_sheet.branch_id_name.name != ""
        ? loginData?.opened_sheet.branch_id_name.name
        : "Branch"

    return (
        <Dropdown
            options={optionsBranches}
            placeholder={placeholderBranch}

            onSelect={(e) => { setBranchName(e.value) }}
            addNewValue={() => { }}
        />
    );
}

export default SelectBranch;
