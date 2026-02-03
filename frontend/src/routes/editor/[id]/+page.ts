interface Props {
  id: string;
}

export const load = ({ url, params }): Props => {
  return {
    id: params.id
  };
};