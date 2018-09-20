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
      'Makroud': 24014,
      'SQLX': 193563,
      'GORM': 89607,
      'GORP': 44242,
      'XORM': 80892,
    },
    'bop': {
      'Makroud': 6584,
      'SQLX': 5031,
      'GORM': 26223,
      'GORP': 8944,
      'XORM': 16268,
    },
    'aop': {
      'Makroud': 64,
      'SQLX': 48,
      'GORM': 418,
      'GORP': 362,
      'XORM': 405,
    },
  },
  'select_subset': {
    'nsop': {
      'Makroud': 25134,
      'SQLX': 113860,
      'GORM': 88795,
      'GORP': 46822,
      'XORM': 75414,
    },
    'bop': {
      'Makroud': 6952,
      'SQLX': 5205,
      'GORM': 27471,
      'GORP': 8944,
      'XORM': 15980,
    },
    'aop': {
      'Makroud': 66,
      'SQLX': 49,
      'GORM': 431,
      'GORP': 362,
      'XORM': 401,
    },
  },
  'select_complex': {
    'nsop': {
      'Makroud': 30429,
      'SQLX': 130819,
      'GORM': 99197,
      'GORP': 43960,
      'XORM': 100789,
    },
    'bop': {
      'Makroud': 8049,
      'SQLX': 5534,
      'GORM': 34009,
      'GORP': 9256,
      'XORM': 17580,
    },
    'aop': {
      'Makroud': 81,
      'SQLX': 57,
      'GORM': 519,
      'GORP': 368,
      'XORM': 445,
    },
  },
  'insert': {
    'nsop': {
      'Makroud': 30669,
      'SQLX': 39335,
      'GORM': 24075,
      'GORP': 5815,
      'XORM': 17758,
    },
    'bop': {
      'Makroud': 6881,
      'SQLX': 2831,
      'GORM': 7928,
      'GORP': 1432,
      'XORM': 5648,
    },
    'aop': {
      'Makroud': 126,
      'SQLX': 49,
      'GORM': 150,
      'GORP': 34,
      'XORM': 127,
    },
  },
  'update': {
    'nsop': {
      'Makroud': 32227,
      'SQLX': 31900,
      'GORM': 48402,
      'GORP': 5136,
      'XORM': 32869,
    },
    'bop': {
      'Makroud': 7410,
      'SQLX': 2463,
      'GORM': 12632,
      'GORP': 1536,
      'XORM': 7640,
    },
    'aop': {
      'Makroud': 123,
      'SQLX': 43,
      'GORM': 287,
      'GORP': 35,
      'XORM': 197,
    },
  },
  'delete': {
    'nsop': {
      'Makroud': 7778,
      'SQLX': 16405,
      'GORM': 15838,
      'GORP': 2084,
      'XORM': 29776,
    },
    'bop': {
      'Makroud': 2120,
      'SQLX': 1215,
      'GORM': 4664,
      'GORP': 352,
      'XORM': 9280,
    },
    'aop': {
      'Makroud': 39,
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
