import { Sheet, TableData, useApp } from "./AppContext";
import Dropdown, { DropdownOption } from "./dropdown";
import { NewNameProps } from "./NewNameModal";

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
        setNewNameModal,
        openedSheet,
    } = useApp();
    const optionsProjects = (loginData?.table_names ?? []).map(item => ({
        value: item.id,
        label: item.name
    }))
    const placeholderProjectas = loginData?.opened_table?.name != ""
        ? loginData?.opened_table?.name
        : "Branch"

    const setData = (data: ProjectData) => {
        setCurrTable(data.Table);
        setCurrSheet(data.Sheet);
        setColumns(data.Sheet.columns);
        setOpenedSheet(data.Sheet);
    }

    return (
        <Dropdown
            options={optionsProjects}
            placeholder={placeholderProjectas}
            onSelect={(e) => getCurrTable(e.value, accessToken ?? "", setData)}
            addNewValue={() => {
                const props: NewNameProps = {
                    currNames: openedSheet?.sheets_id_names ?? [],
                    assignNewName: (name) => postTable(name, accessToken ?? "", setData),
                }
                setNewNameModal(props)
            }}
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
                throw new Error("Could not retrieve project");
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

const postTable = (name: string, token: string, setData: Function) => {
    const url = `/create_project`;
    const nameParam: { Name: string } = { Name: name }
    fetch(url, {
        method: "POST",
        headers: {
            'Authorization': `Bearer ${token}`
        },
        credentials: "include",
        body: JSON.stringify(nameParam)
    })
        .then(response => {
            if (response.status !== 201) {
                throw new Error("Could not create project");
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
