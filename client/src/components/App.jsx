import React, { useState, useEffect } from 'react';
import axios from 'axios';
import DeleteOutlineRoundedIcon from '@material-ui/icons/DeleteOutlineRounded';
import MenuBookRoundedIcon from '@material-ui/icons/MenuBookRounded';
import AddRoundedIcon from '@material-ui/icons/AddRounded';
import Fab from '@material-ui/core/Fab';

const endpoint = "https://gowiki-api.herokuapp.com";

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

// A function to render spaces between the create-area and the articles
function Spacing() {
    return (
        <div className="spacing">
            <br></br>
            <br></br>
        </div>
    );
}

// Render the articles after they are created
function Article(props) {
    return (
        <div className="article">
            <h1>{props.title}</h1>
            <p>{props.body}</p>
            <img src={props.image} alt="" width="240" height="180"/>
            <p>{props.author}</p>
            <button onClick={() => props.onDelete(props.id)}>
                <DeleteOutlineRoundedIcon />
            </button>
        </div>
    );
}

// Render the CreateArea, where you can enter the contents of the article
function CreateArea(props) {
    const [titleText, setTitleText] = useState("");
    const [bodyText, setBodyText] = useState("");
    const [imageURI, setImageURI] = useState("");
    const [authorName, setAuthorName] = useState("");

    return (
        <div>
            <form className="create-article">
                <input 
                    value={titleText}
                    onChange={(event) => setTitleText(event.target.value)}
                    name="title"
                    placeholder="Title"
                    autoComplete="off"
                    required
                />
                <textarea 
                    value={bodyText}
                    onChange={(event) => setBodyText(event.target.value)}
                    name="body"
                    placeholder="Type something here"
                    required
                />
                <textarea 
                    value={imageURI}
                    onChange={(event) => setImageURI(event.target.value)}
                    name="image"
                    placeholder="Optional Image URL"
                />
                <textarea 
                    value={authorName}
                    onChange={(event) => setAuthorName(event.target.value)}
                    name="author"
                    placeholder="Author"
                    required
                />
                <Fab color="secondary" aria-label="add" onClick={(event) => {
                    props.onAdd({ title: titleText, body: bodyText, image: imageURI, author: authorName });
                    setTitleText("");
                    setBodyText("");
                    setImageURI("");
                    setAuthorName("");
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
    function onDelete(id) {
        axios.delete(`${endpoint}/api/delete/${id}`).then(() => {
            setArticles((prev) => prev.filter(note._id != id))
        })
    }

    useEffect(() => {
        axios.get(endpoint + "/api/read").then((res) => {    
            setArticles(res.data);
        });
    }, []);

return (
    <div>
        <Header />
        <Spacing />
        <CreateArea onAdd={createArticle} />
        <Spacing /><Spacing />
        {articles.map((article) => (
            <Article 
                key={article._id}
                id={article._id}
                title={article.title}
                body={article.body}
                image={article.image}
                author={article.author}
                onDelete={onDelete}
            />
        ))}
    </div>
    );
}

export default App;
