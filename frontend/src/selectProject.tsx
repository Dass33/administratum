import { Sheet, TableData, useApp } from "./AppContext";
import Dropdown, { DropdownOption } from "./dropdown";

type ProjectData = {
    Table: TableData,
    Sheet: Sheet,
}

function SelectProject() {
    const {
        loginData,
        accessToken,
        setCurrTable,
        setCurrSheet,
        setColumns,
        setOpenedSheet,
    } = useApp();
    const optionsProjects = (loginData?.table_names ?? []).map(item => ({
        value: item.id,
        label: item.name
    }))
    const placeholderProjectas = loginData?.opened_table?.name != ""
        ? loginData?.opened_table?.name
        : "Branch"

    const setProject = (e: DropdownOption) => {
        getCurrTable(e.value, accessToken ?? "", (data: ProjectData) => {
            setCurrTable(data.Table);
            setCurrSheet(data.Sheet);
            setColumns(data.Sheet.columns);
            setOpenedSheet(data.Sheet);
        })
    }

    return (
        <Dropdown
            options={optionsProjects}
            placeholder={placeholderProjectas}
            onSelect={setProject}
            addNewValue={() => { }}
        />
    );
}

const getCurrTable = (table_id: string, token: string, setData: Function) => {
    const url = `/get_project/${table_id}`;

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
        .then(result => {
            setData(result);
        })
        .catch(err => {
            console.error(err);
        });
};

export default SelectProject;
