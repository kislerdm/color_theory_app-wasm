"""Training script"""
import logging
import pathlib
from typing import Any

import pandas
import xgboost
from sklearn.model_selection import train_test_split


def train_model(x_train: Any, y_train: Any, seed: int) -> xgboost.Booster:
    """ Trains the model.

    Args:
     x_train: Data to train the model.
     y_train: Labels to train the model.
     seed: Seed to randomise split.

    Returns:
        Model.
    """
    params = {
        "objective": 'binary:logistic',
        "learning_rate": 0.5,
        "max_depth": 3,
        "n_jobs": 4,
        "subsample": 0.8,
        "random_state": seed,
        "eval_metric": "error",
    }

    num_round = 10

    data_train = xgboost.DMatrix(x_train, label=y_train, silent=True)
    return xgboost.train(params, data_train, num_round, verbose_eval=True)


def main(path_data: str, path_model: str, seed: int) -> None:
    """Main routine

    Args:
        path_data: Path to data csv file.
        path_model: Path to save the model.
        seed: Seed to randomise split.
    """
    data = pandas.read_csv(path_data) \
        .sample(frac=1, random_state=seed) \
        .reset_index(drop=True)

    x = data.drop('is_warm', axis=1)
    y = data.get('is_warm')

    x_train, _, y_train, _ = train_test_split(x, y, test_size=0.2, random_state=seed)

    model = train_model(x_train, y_train, seed)

    model.dump_model(path_model, dump_format="json")


if __name__ == "__main__":
    base = pathlib.Path(__file__).absolute().parent

    path_input = f"{base}/warm_cold_colors.csv"
    path_model = f"{base}/model.json"
    seed = 2019

    log = logging.Logger(__file__)
    try:
        main(path_input, path_model, seed)
    except Exception as ex:
        log.error(ex)
