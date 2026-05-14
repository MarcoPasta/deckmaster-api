import { useSession } from "@clerk/vue";
import { ref } from "vue";
import { api, ApiError } from "../services/api";

export function useApi() {
  const { session } = useSession();
  const loading = ref(false);
  const error = ref<ApiError | null>(null);

  async function getToken() {
    return (await session.value?.getToken()) ?? null;
  }

  async function get<T>(path: string, params?: Record<string, string>) {
    loading.value = true;
    error.value = null;
    try {
      const token = await getToken();
      return await api.get<T>(path, token, { params });
    } catch (e) {
      error.value = e instanceof ApiError ? e : null;
      throw e;
    } finally {
      loading.value = false;
    }
  }

  async function post<T>(path: string, body?: unknown) {
    loading.value = true;
    error.value = null;
    try {
      const token = await getToken();
      return await api.post<T>(path, token, { body });
    } catch (e) {
      error.value = e instanceof ApiError ? e : null;
      throw e;
    } finally {
      loading.value = false;
    }
  }

  async function put<T>(path: string, body?: unknown) {
    loading.value = true;
    error.value = null;
    try {
      const token = await getToken();
      return await api.put<T>(path, token, { body });
    } catch (e) {
      error.value = e instanceof ApiError ? e : null;
      throw e;
    } finally {
      loading.value = false;
    }
  }

  async function del<T>(path: string) {
    loading.value = true;
    error.value = null;
    try {
      const token = await getToken();
      return await api.delete<T>(path, token);
    } catch (e) {
      error.value = e instanceof ApiError ? e : null;
      throw e;
    } finally {
      loading.value = false;
    }
  }

  async function postForm<T>(path: string, formData: FormData) {
    loading.value = true;
    error.value = null;
    try {
      const token = await getToken();
      return await api.postForm<T>(path, token, formData);
    } catch (e) {
      error.value = e instanceof ApiError ? e : null;
      throw e;
    } finally {
      loading.value = false;
    }
  }

  async function getBlob(path: string, params?: Record<string, string>) {
    loading.value = true;
    error.value = null;
    try {
      const token = await getToken();
      return await api.getBlob(path, token, { params });
    } catch (e) {
      error.value = e instanceof ApiError ? e : null;
      throw e;
    } finally {
      loading.value = false;
    }
  }

  return { get, post, put, del, postForm, getBlob, loading, error };
}
