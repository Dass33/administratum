import expand from "./assets/expand.svg"

function GameView() {

    return (
        <div className="flex flex-row">
            <div className="embedded-container">
                <iframe className="rounded-lg"
                    src="https://dass33.github.io/guess_game/"
                    width="370"
                    height="750"
                    title="Embedded Website"
                />
            </div>
            <img
                className="mx-4 size-8 hover:scale-125 transition-transform duration-200"
                src={expand}
                alt="dropdown arrow"
            />
        </div>
    )
}

export default GameView
