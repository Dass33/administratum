import expand from "./assets/expand.svg"

function GameView() {

    return (
        <div className="flex items-center">
            <div className="flex flex-row">
                <button className="ml-4 size-8 hover:scale-125 transition-transform duration-100" >
                    <img
                        src={expand}
                        alt="dropdown arrow"
                    />
                </button>
                <div className="embedded-container">
                    <iframe className="rounded-lg"
                        src="https://dass33.github.io/guess_game/"
                        width="370"
                        height="750"
                        title="Embedded Website"
                    />
                </div>
            </div>
        </div>
    )
}

export default GameView
