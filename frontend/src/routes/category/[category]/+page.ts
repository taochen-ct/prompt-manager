interface Props {
  category: string;
}

export const load = ({ url, params }): Props => {
  return {
    category: params.category
  };
};