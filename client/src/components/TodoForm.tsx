import { BsPlusCircle } from "react-icons/bs";
import {useState} from "react";
import {BASE_URL, Todo} from "../App.tsx";
import * as React from "react";
import {CircularProgress} from "@mui/material";

export default function TodoForm({addTodo}:{addTodo:(newTodo:Todo)=>void}) {
    const [todoBody,setTodoBody] = useState<string>("");
    const [isCreating,setIsCreating] = useState<Boolean>(false);

    async function createTodo(event:React.FormEvent) {
        setIsCreating(true);
        event.preventDefault();
        try {
            const res = await fetch(BASE_URL + "/todos", {
                method:"POST",
                headers: {
                    "content-type": "application/json"
                },
                body:JSON.stringify({body:todoBody})
            });
            const data = await res.json();
            if(!res.ok) {
                throw new Error(data.error || "createTodo failed");
            }
            addTodo(data);
            setTodoBody("");
            return data;
        }
        catch (e) {
            console.log(e);
        }
        finally {
            setIsCreating(false);
        }
    }

    return (
        <div id="main-div">
            <form onSubmit={createTodo}>
                <input name="todoBody" value={todoBody} placeholder="Add task" autoFocus onChange={(e)=>setTodoBody(e.target.value)}></input>
                <button id="create-button" type="submit">
                    {isCreating ? <CircularProgress size={13}/> : <BsPlusCircle />}
                </button>
            </form>
        </div>
    )
}