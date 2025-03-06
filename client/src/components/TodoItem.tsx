import { FaRegCheckCircle } from "react-icons/fa";
import { FaRegTrashCan } from "react-icons/fa6";
import {BASE_URL, Todo} from "../App.tsx";
import {useState} from "react";
import {CircularProgress} from "@mui/material";

export default function TodoItem({todo, changeTodo, removeTodo}:{todo:Todo, changeTodo:(updatedTodo:Todo)=>void, removeTodo:(deletedTodo:Todo)=>void}) {
    const [isUpdating, setIsUpdating] = useState<Boolean>(false);
    const [isDeleting, setIsDeleting] = useState<Boolean>(false);

    async function updateTodo() {
        setIsUpdating(true);
        try {
            if (todo.completed) {
                return alert("Todo already completed!");
            }
            const res = await fetch(BASE_URL + `/todos/${todo._id}`, {
                method: "PATCH"
            });
            const data = await res.json();
            if (!res.ok) {
                throw new Error(data.error || "updateTodo failed");
            }
            todo.completed=true;
            changeTodo(todo);
            return data;
        }
        catch (e) {
            console.log(e);
        }
        finally {
            setIsUpdating(false);
        }
    }

    async function deleteTodo() {
        setIsDeleting(true);
        try {
            const res = await fetch(BASE_URL + `/todos/${todo._id}`, {
                method:"DELETE"
            });
            const data = await res.json();
            if (!res.ok) {
                throw new Error(data.error || "deleteTodo failed")
            }
            removeTodo(todo);
            return data;
        }
        catch (e) {
            console.log(e);
        }
        finally {
            setIsDeleting(false);
        }
    }

    return (
        <div className="todoItem-outer-div">
            <div className="todoItem-inner-div1">{todo.body}</div>
            <div className="todoItem-inner-div2"
                 style={{backgroundColor: todo.completed ? "green" : "orange"}}
            >
                {todo.completed ? "Done" : "In progress"}
            </div>
            <button className="update-button" type="submit" onClick={updateTodo}>
                {isUpdating ? <CircularProgress size={13}/> : <FaRegCheckCircle />}
            </button>
            <button className="delete-button" type="submit" onClick={deleteTodo}>
                {isDeleting ? <CircularProgress size={13}/> : <FaRegTrashCan />}
            </button>
        </div>
    )
}