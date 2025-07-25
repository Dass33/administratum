import BottomBar from "./bottomBar";
import GameView from "./gameView"
import SaveButton from "./saveButton";
import SelectBranch from "./selectBranch";
import SelectProject from "./selectProject";
import Permissions from "./permissions";
import Table from "./table";
import { useApp } from "./AppContext";
import CellModal from "./cellModal";

GameView
function App() {
    const { showCellModal } = useApp()

    const sampleData = [
        { name: "John", age: 30, city: "New York", active: true },
        { name: "Jane", age: 25, city: "Los Angeles", department: "Engineering" },
        { name: "Bob", age: 35, active: false, salary: 75000 },
        { name: "John", age: 30, city: "New York", active: true },
        { name: "Jane", age: 25, city: "Los Angeles", department: "Engineering" },
        { name: "Bob", age: 35, active: false, salary: 75000 },
        { name: "John", age: 30, city: "New York", active: true },
        { name: "Jane", age: 25, city: "Los Angeles", department: "Engineering" },
        { name: "Bob", age: 35, active: false, salary: 75000 },
        { name: "John", age: 30, city: "New York", active: true },
        { name: "Jane", age: 25, city: "Los Angeles", department: "Engineering" },
        { name: "Bob", age: 35, active: false, salary: 75000, questions: "this is a quesetions?", config: "hhhellaaaaaaaakfjklsadfaellaaaaaaaakfjklsadfahellaaaaaaaakfjklsadfaellaaaaaaaakfjklsadfa" },
    ];
    return (
        <div className='bg-figma-white h-screen pl-10 pr-5 flex flex-row xl:gap-16'>
            {showCellModal && <CellModal />}
            <div className="flex flex-col flex-grow my-10">
                <div className="flex flex-row gap-4">
                    <SaveButton />

                    <div className="flex items-center text-lg">
                        <SelectProject />
                        <span className="mx-2">/</span>
                        <SelectBranch />
                    </div>
                    <Permissions />
                </div>
                <div className="mt-10 grow">
                    <Table data={sampleData} />
                </div>
                <div className="bottom-1 absolute my-2">
                    <BottomBar />
                </div>
            </div>
            <GameView />
        </div>
    )
}

export default App
