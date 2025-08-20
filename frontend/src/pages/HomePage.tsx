import '../styles/index.scss'
import '../styles/homepage.scss'

export default function HomePage() {

    return (
        <div className="homepage__container">
            <div className="homepage__title">
                <h1 className="homepage__title--main">tomomo</h1>
                <div className="homepage__title--sub">
                    <div className="title--sub--first">todo</div>-
                    <div className="title--sub--second">mood</div>-
                    <div className="title--sub--third">money</div>
                </div>
            </div>
        </div>
    );
}
