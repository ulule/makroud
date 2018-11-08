#!/usr/bin/python3

import plotly.offline as py
import plotly.io as pio
import plotly.graph_objs as go

import os

def make_graph(filename, title, ytitle, benchmark):

  data = [
    go.Bar(
      x = list(benchmark.keys()),
      y = list(benchmark.values()),
      marker = {
        'color': [
          'rgb(49,171,95)',
          'rgb(49, 110, 171)',
          'rgb(212, 109, 57)',
          'rgb(148, 62, 154)',
          'rgb(54, 176, 165)'
        ],
      }
    ),
  ]

  layout = go.Layout(
    title = title,
    width = 400,
    height = 400,
    xaxis = {
      'title': 'ORM',
    },
    yaxis = {
      'title': ytitle,
    }
  )

  fig = go.Figure(data=data, layout=layout)
  pio.write_image(fig, filename)


title_table = {'nsop': 'Speed', 'bop': 'Memory', 'aop': 'Allocations'}
yaxis_table = {'nsop': 'ns/op', 'bop': 'B/op', 'aop': 'allocs/op'}

benchmark_table = {
  'select_all': {
    'nsop': {
      'Makroud': 15292,
      'SQLX': 162435,
      'GORM': 62021,
      'GORP': 29564,
      'XORM': 56527,
    },
    'bop': {
      'Makroud': 6288,
      'SQLX': 4767,
      'GORM': 26367,
      'GORP': 8944,
      'XORM': 16332,
    },
    'aop': {
      'Makroud': 62,
      'SQLX': 45,
      'GORM': 418,
      'GORP': 362,
      'XORM': 405,
    },
  },
  'select_subset': {
    'nsop': {
      'Makroud': 16003,
      'SQLX': 164521,
      'GORM': 64113,
      'GORP': 29786,
      'XORM': 53138,
    },
    'bop': {
      'Makroud': 6480,
      'SQLX': 4767,
      'GORM': 27615,
      'GORP': 8944,
      'XORM': 16044,
    },
    'aop': {
      'Makroud': 64,
      'SQLX': 45,
      'GORM': 431,
      'GORP': 362,
      'XORM': 401,
    },
  },
  'select_complex': {
    'nsop': {
      'Makroud': 17167,
      'SQLX': 331660,
      'GORM': 74538,
      'GORP': 30402,
      'XORM': 58262,
    },
    'bop': {
      'Makroud': 6937,
      'SQLX': 4887,
      'GORM': 34215,
      'GORP': 9256,
      'XORM': 17644,
    },
    'aop': {
      'Makroud': 74,
      'SQLX': 48,
      'GORM': 519,
      'GORP': 368,
      'XORM': 445,
    },
  },
  'insert': {
    'nsop': {
      'Makroud': 15536,
      'SQLX': 34789,
      'GORM': 16785,
      'GORP': 4350,
      'XORM': 12657,
    },
    'bop': {
      'Makroud': 5673,
      'SQLX': 2831,
      'GORM': 7184,
      'GORP': 1592,
      'XORM': 5872,
    },
    'aop': {
      'Makroud': 109,
      'SQLX': 49,
      'GORM': 146,
      'GORP': 37,
      'XORM': 130,
    },
  },
  'update': {
    'nsop': {
      'Makroud': 16924,
      'SQLX': 28036,
      'GORM': 33415,
      'GORP': 3577,
      'XORM': 18940,
    },
    'bop': {
      'Makroud': 5849,
      'SQLX': 2463,
      'GORM': 12712,
      'GORP': 1536,
      'XORM': 7704,
    },
    'aop': {
      'Makroud': 104,
      'SQLX': 43,
      'GORM': 287,
      'GORP': 35,
      'XORM': 197,
    },
  },
  'delete': {
    'nsop': {
      'Makroud': 3698,
      'SQLX': 12555,
      'GORM': 11101,
      'GORP': 1505,
      'XORM': 20893,
    },
    'bop': {
      'Makroud': 1392,
      'SQLX': 1215,
      'GORM': 4728,
      'GORP': 352,
      'XORM': 9344,
    },
    'aop': {
      'Makroud': 32,
      'SQLX': 22,
      'GORM': 95,
      'GORP': 13,
      'XORM': 202,
    },
  },
}

key_table = {
  'SelectAll': 'select_all',
  'SelectSubset': 'select_subset',
  'SelectComplex': 'select_complex',
  'Insert': 'insert',
  'Update': 'update',
  'Delete': 'delete',
}

if not os.path.exists('images'):
  os.mkdir('images')

for kind in ['SelectAll', 'SelectSubset', 'SelectComplex', 'Insert', 'Update', 'Delete']:
  for bench in ['nsop', 'bop', 'aop']:
    make_graph(
      'images/%s_%s.png' % (key_table[kind], bench),
      '%s %s' % (kind, title_table[bench]),
      yaxis_table[bench],
      benchmark_table[key_table[kind]][bench],
    )
