import { Container, Group } from "@mantine/core";
import useSWR from "swr";

import { serverURL } from "src/config/constants";
import { fetcher } from "src/api/api";
import { ITodo } from "src/model/interfaces";
import { Todo } from "src/components/todo";
import * as API from "src/api/api";
import { Add } from "src/components/add";

const Home = () => {
  const getTodos = () => {
    const { data, error, isLoading, mutate } = useSWR(
      `${serverURL}/todos`,
      fetcher
    );

    return {
      todos: data,
      isLoading,
      mutate,
      isError: error,
    };
  };

  const { todos, isLoading, isError, mutate } = getTodos();

  const updateTodo = async (id: string, checked: boolean) => {
    const response = await API.updateTodo(checked, id);
    if (response.status != 201 && response.status != 200) {
      alert(response?.data?.message);
    } else {
      mutate(`${serverURL}/todos`);
    }
  };

  const createTodo = async (title: string) => {
    const response = await API.createTodo(title);
    if (response.status != 201 && response.status != 200) {
      alert(response?.data?.message);
    } else {
      mutate(`${serverURL}/todos`);
    }
  };

  return isLoading ? (
    <h1>Is Loading...</h1>
  ) : (
    <div id="todos">
      <Container size={"md"}>
        <h2 style={{ textAlign: "left", marginTop: 50 }}>Todos</h2>

        <Group>
          <Add onAdd={(title: string) => createTodo(title)} />
          {todos &&
            Array.isArray(todos) &&
            todos?.map((todo: ITodo, index: number) => (
              <Todo
                key={index}
                title={todo.title}
                description={todo.date}
                checked={todo.done}
                onChange={(checked: boolean) => {
                  updateTodo(todo.id, checked);
                }}
              />
            ))}
        </Group>
      </Container>
    </div>
  );
};
export default Home;
