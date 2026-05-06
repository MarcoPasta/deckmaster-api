import { ref } from "vue";

export type TcgType = "pokemon" | "naruto" | "digimon";

const selectedTcg = ref<TcgType>(
  (sessionStorage.getItem("selectedTcg") as TcgType) || "pokemon",
);

export function useTcgState() {
  function setTcg(tcg: TcgType) {
    selectedTcg.value = tcg;
    sessionStorage.setItem("selectedTcg", tcg);
  }

  return {
    selectedTcg,
    setTcg,
  };
}
