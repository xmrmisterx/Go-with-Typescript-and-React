import Header from "./components/Header.tsx";
import TodoForm from "./components/TodoForm.tsx";
import TodoList from "./components/TodoList.tsx";
import {useState} from "react";

export type Todo = {
    _id:number,
    completed:boolean,
    body:string
}

export const BASE_URL="http://localhost:5005/api";

function App() {
    const [todos,setTodos] = useState<Todo[]>([]);

    function addTodo(newTodo:Todo) {
        setTodos(prevTodos=>{
            return [...prevTodos,newTodo];
        });
    }

    return (
      <>
        <Header />
        <TodoForm addTodo={addTodo}/>
        <TodoList todos={todos} setTodos={setTodos}/>
      </>
  )
}

export default App
