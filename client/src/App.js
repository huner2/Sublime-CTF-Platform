import React from 'react';
//import logo from './logo.svg';
import './App.css';
import NavBar from './Components/NavBar';
import { Route } from 'react-router-dom';
import Home from './Components/Home';

class App extends React.Component {
  constructor(props){
    super(props);
    this.state = {upath: props.location.pathname};
    console.log(this.state);
  }

  render() {
    return (
      <div className="App">
        <NavBar path={this.state.upath}/>
        {/*<Route path="/about" component={About} />}*/}
        <Route path="/" component={Home} />
      </div>
    );
  }
}

export default App;
