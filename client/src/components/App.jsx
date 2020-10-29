import React, { useState, useEffect } from 'react';
import axios from 'axios';
import DeleteOutlineIcon from '@material-ui/icons/DeleteOutline'
import MenuBookRoundedIcon from '@material-ui/icons/MenuBookRounded';
import AddRoundedIcon from '@material-ui/icons/AddRounded';
import Fab from '@material-ui/core/Fab';

const endpoint = "http://localhost:8080";

// A function to render the Header of this webpage
function Header() {
    return (
        <header className="header">
            <div class="jumbotron jumbotron-fluid text-center">
                <h1 class="display-4"><MenuBookRoundedIcon fontSize="large"/> GoWiki</h1>
                <hr class="rule"></hr>
                <h5>A website where you can create, view,</h5>
                <h5>and edit science articles</h5>
            </div>
        </header>
    );
}

// Render the articles after they are created
function Article(props) {
    return (
        <div className="article">
            <h1>{props.title}</h1>
            <p className="lead">{props.body}</p>
            <button onClick={() => props.onDelete(props.title)}>
                <DeleteOutlineIcon />
            </button>
        </div>
    );
}

// Render the CreateArea, where you can enter the contents of the article
function CreateArea(props) {
    const [titleText, setTitleText] = useState("");
    const [bodyText, setBodyText] = useState("");

    return (
        <div>
            <form className="create-article">
                <input 
                    value={titleText}
                    onChange={(event) => setTitleText(event.target.value)}
                    name="title"
                    placeholder="Title"
                    autoComplete="off"
                />
                <textarea 
                    value={bodyText}
                    onChange={(event) => setBodyText(event.target.value)}
                    name="body"
                    placeholder="Type something here"
                />
                <Fab color="secondary" aria-label="add" onClick={(event) => {
                    props.onAdd({ title: titleText, body: bodyText });
                    setTitleText("");
                    setBodyText("");
                    event.preventDefault();
                }}>
                    <AddRoundedIcon />
                </Fab>
            </form>
        </div>
    );
}

function App() {
    // Make use of states, in this case a list of articles
    const [articles, setArticles] = useState([]);

    // A function to create an article using a POST request
    function createArticle(article) {
        axios.post(endpoint + "/api/create", article);
        setArticles((prev) => [...prev, article])
    }

    // A function to delete an article using a DELETE request
    function onDelete(title) {
        axios.delete(`${endpoint}/delete/${title}`)
    }

    useEffect(() => {
        axios.get(endpoint + "/api/read").then((res) => {    
            setArticles(res.data);
        });
    }, []);

return (
    <div>
        <Header />
        <CreateArea onAdd={createArticle} />
        {articles.map((article) => (
            <Article 
                key={article._id}
                id={article._id}
                title={article.title}
                body={article.body}
                onDelete={onDelete}
            />
        ))}
    </div>
    );
}

export default App;