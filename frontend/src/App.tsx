import BottomBar from "./bottomBar";
import GameView from "./gameView"
import SaveButton from "./saveButton";
import SelectBranch from "./selectBranch";
import SelectProject from "./selectProject";
import Permissions from "./permissions";
import Table from "./table";
import { useApp } from "./AppContext";
import CellModal from "./cellModal";
import ColModal from "./colModal";

GameView
function App() {
    const { cellModal, colModal } = useApp()

    return (
        <div className='bg-figma-white h-screen pl-10 xl:pr-5 flex flex-row'>
            {cellModal && <CellModal />}
            {colModal && <ColModal />}
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
                    <Table />
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
