import TodoItem from "./TodoItem.tsx";
import {SetStateAction, useEffect, useState} from "react";
import {BASE_URL, Todo} from "../App.tsx";
import * as React from "react";
import {CircularProgress} from "@mui/material";

export default function TodoList({todos,setTodos}:{todos:Todo[],setTodos:React.Dispatch<SetStateAction<Todo[]>>}) {
	const [isLoading,setIsLoading] = useState<Boolean>(false);

	useEffect(()=>{
		setIsLoading(true);
		fetch(BASE_URL + "/todos")
			.then(res=>res.json())
			.then(data=>{
				// console.log("data:", data);
				data ? setTodos(data) : setTodos([]);
			})
			.finally(()=>{
				setIsLoading(false);
			});
	},[]);

	function changeTodo(updatedTodo:Todo) {
		setTodos(prevTodos=>
			prevTodos.map(todo=>todo._id===updatedTodo._id ? updatedTodo : todo)
		)
	}

	function removeTodo(deletedTodo:Todo) {
		setTodos(prevTodos =>
			prevTodos.filter(todo=> todo._id!==deletedTodo._id)
		)
	}

    return (
        <div id="todoList-div">
            <h2>Today's Tasks</h2>
			{isLoading && <CircularProgress />}
			{!isLoading && todos.length===0 && <h3>ALL TASKS COMPLETE!</h3>}

			{todos.map(todo=> (
				<TodoItem key={todo._id} todo={todo} changeTodo={changeTodo} removeTodo={removeTodo}/>
			))}
        </div>

    )
}