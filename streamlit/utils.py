"""Utils module."""
import pandas as pd


def safe_to_numeric(series):
    """Cast numeric serie."""
    try:
        return pd.to_numeric(series)
    except (ValueError, TypeError):
        return series
