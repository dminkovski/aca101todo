import * as config from "src/config/constants";
import axios from "axios";

export async function updateTodo(checked: boolean, id: string) {
  const data = {
    done: checked,
  };
  try {
    const response = await axios.post(`${config.serverURL}/todos/${id}`, data, {
      headers: { "Content-Type": "application/json" },
    });
    console.log("response", response);
    return response;
  } catch (error: any) {
    console.log("error", error);
    if (error?.response) {
      return error.response;
    }
    return error;
  }
}

export async function createTodo(title: string) {
  const data = {
    title: title,
  };
  try {
    const response = await axios.post(`${config.serverURL}/todos`, data, {
      headers: { "Content-Type": "application/json" },
    });
    console.log("response", response);
    return response;
  } catch (error: any) {
    console.log("error", error);
    if (error?.response) {
      return error.response;
    }
    return error;
  }
}

export const fetcher = (url: string) => axios.get(url).then((res) => res.data);
