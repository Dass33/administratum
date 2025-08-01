import { useApp } from "./AppContext";
import Dropdown from "./dropdown";

function SelectProject() {
    const { setProjectName, loginData } = useApp();
    const optionsProjects = (loginData?.table_names ?? []).map(item => ({
        value: item.id,
        label: item.name
    }))
    const placeholderProjectas = loginData?.opened_table?.name != ""
        ? loginData?.opened_table?.name
        : "Branch"

    return (
        <Dropdown
            options={optionsProjects}
            placeholder={placeholderProjectas}
            onSelect={(e) => { setProjectName(e.value) }}
            addNewValue={() => { }}
        />
    );
}

export default SelectProject;
