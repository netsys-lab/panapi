#! /bin/python3

import time
import plotly.express as px
import pandas as pd
import dash
import dash_core_components as dcc
import dash_html_components as html
from dash.dependencies import Input, Output

app = dash.Dash()
app.layout = html.Div([
    html.Div(id='live-bandwidth'),
    dcc.Graph(id='live-figure'),
    dcc.Interval(
        id='interval-component',
        interval=500, # in milliseconds
        n_intervals=0
    )
])

@app.callback(Output('live-bandwidth', 'children'),
              Input('interval-component', 'n_intervals'))
def update_metrics(n):
    bandwidth = pd.read_csv("spate/pid.csv")['Mibps'].mean()
    return [html.Span('Bandwidth: {0:f} Mib/s'.format(bandwidth))]

@app.callback(Output('live-figure', 'figure'),
              Input('interval-component', 'n_intervals'))
def updage_figure(n):
    df = pd.read_csv("spate/pid.csv")
    fig = px.line(df, y = 'Mibps')

    return fig

app.run_server()
