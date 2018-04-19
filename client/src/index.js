import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import App from './App';
import registerServiceWorker from './registerServiceWorker';

import { BrowserRouter as Router, Route } from 'react-router-dom';

const Root = () => (
  <Router>
    <Route path="/" component={App} />
  </Router>
)

ReactDOM.render(<Root />, document.getElementById('root'));
registerServiceWorker();
