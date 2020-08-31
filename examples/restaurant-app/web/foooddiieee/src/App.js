import React from 'react';
import './App.css';
import TopNav from './components/TopNav'
import Intro from './components/Intro'
import Food from './components/Food'
import Menu from './components/Menu'
import Info from './components/Info'
import Manage from './components/Manage'
import Footer from './components/Footer'
import {AppContextProvider} from './components/AppContext'

function App() {

  return (
    <AppContextProvider>
      <div className="App">
        <TopNav/ >
        <Intro/ >
        <Menu/ >
        <Food/ >
        <Info/ >
        <Footer/ >
      </div>
    </AppContextProvider>
  );
}

export default App;
