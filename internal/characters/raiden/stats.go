package raiden

var (
	attack = [][][]float64{
		//1
		{
			{
				0.3965,
				0.4287,
				0.461,
				0.5071,
				0.5394,
				0.5763,
				0.627,
				0.6777,
				0.7284,
				0.7837,
				0.8471,
				0.9216,
				0.9962,
				1.0707,
				1.152,
			},
		},
		//2
		{
			{
				0.3973,
				0.4297,
				0.462,
				0.5082,
				0.5405,
				0.5775,
				0.6283,
				0.6791,
				0.73,
				0.7854,
				0.8489,
				0.9236,
				0.9983,
				1.073,
				1.1545,
			},
		},
		//3
		{
			{
				0.4988,
				0.5394,
				0.58,
				0.638,
				0.6786,
				0.725,
				0.7888,
				0.8526,
				0.9164,
				0.986,
				1.0658,
				1.1595,
				1.2533,
				1.3471,
				1.4494,
			},
		},
		//4
		{
			{
				0.2898,
				0.3134,
				0.337,
				0.3707,
				0.3943,
				0.4213,
				0.4583,
				0.4954,
				0.5325,
				0.5729,
				0.6192,
				0.6737,
				0.7282,
				0.7827,
				0.8422,
			},
			{
				0.2898,
				0.3134,
				0.337,
				0.3707,
				0.3943,
				0.4213,
				0.4583,
				0.4954,
				0.5325,
				0.5729,
				0.6192,
				0.6737,
				0.7282,
				0.7827,
				0.8422,
			},
		},
		//5
		{
			{
				0.6545,
				0.7077,
				0.761,
				0.8371,
				0.8904,
				0.9513,
				1.035,
				1.1187,
				1.2024,
				1.2937,
				1.3983,
				1.5214,
				1.6444,
				1.7675,
				1.9017,
			},
		},
	}
	attackB = [][][]float64{
		//1
		{
			{
				0.4474,
				0.4779,
				0.5084,
				0.5491,
				0.5796,
				0.6151,
				0.6609,
				0.7066,
				0.7524,
				0.7982,
				0.8439,
				0.8897,
				0.9354,
				0.9812,
				1.0269,
			},
		},
		//2
		{
			{
				0.4396,
				0.4695,
				0.4995,
				0.5395,
				0.5694,
				0.6044,
				0.6494,
				0.6943,
				0.7393,
				0.7842,
				0.8292,
				0.8741,
				0.9191,
				0.964,
				1.009,
			},
		},
		//3
		{
			{
				0.5382,
				0.5749,
				0.6116,
				0.6605,
				0.6972,
				0.74,
				0.7951,
				0.8501,
				0.9052,
				0.9602,
				1.0153,
				1.0703,
				1.1254,
				1.1804,
				1.2355,
			},
		},
		//4
		{
			{
				0.3089,
				0.3299,
				0.351,
				0.3791,
				0.4001,
				0.4247,
				0.4563,
				0.4879,
				0.5195,
				0.5511,
				0.5827,
				0.6143,
				0.6458,
				0.6774,
				0.709,
			},
			{
				0.3098,
				0.3309,
				0.352,
				0.3802,
				0.4013,
				0.4259,
				0.4576,
				0.4893,
				0.521,
				0.5526,
				0.5843,
				0.616,
				0.6477,
				0.6794,
				0.711,
			},
		},
		//5
		{
			{
				0.7394,
				0.7899,
				0.8403,
				0.9075,
				0.9579,
				1.0167,
				1.0924,
				1.168,
				1.2436,
				1.3192,
				1.3948,
				1.4705,
				1.5461,
				1.6217,
				1.6973,
			},
		},
	}
	charge = []float64{
		0.9959,
		1.0769,
		1.158,
		1.2738,
		1.3549,
		1.4475,
		1.5749,
		1.7023,
		1.8296,
		1.9686,
		2.1278,
		2.3151,
		2.5023,
		2.6896,
		2.8938,
	}
	chargeSword = [][]float64{
		{
			0.616,
			0.658,
			0.7,
			0.756,
			0.798,
			0.847,
			0.91,
			0.973,
			1.036,
			1.099,
			1.162,
			1.225,
			1.288,
			1.351,
			1.414,
		},
		{
			0.7436,
			0.7943,
			0.845,
			0.9126,
			0.9633,
			1.0225,
			1.0985,
			1.1746,
			1.2506,
			1.3267,
			1.4027,
			1.4788,
			1.5548,
			1.6309,
			1.7069,
		},
	}
	skill = []float64{
		1.172,
		1.2599,
		1.3478,
		1.465,
		1.5529,
		1.6408,
		1.758,
		1.8752,
		1.9924,
		2.1096,
		2.2268,
		2.344,
		2.4905,
		2.637,
		2.7835,
	}
	skillTick = []float64{
		0.42,
		0.4515,
		0.483,
		0.525,
		0.5565,
		0.588,
		0.63,
		0.672,
		0.714,
		0.756,
		0.798,
		0.84,
		0.8925,
		0.945,
		0.9975,
	}
	skillBurstBonus = []float64{
		0.0022,
		0.0023,
		0.0024,
		0.0025,
		0.0026,
		0.0027,
		0.0028,
		0.0029,
		0.003,
		0.003,
		0.003,
		0.003,
		0.003,
		0.003,
		0.003,
	}
	burstBase = []float64{
		4.008,
		4.3086,
		4.6092,
		5.01,
		5.3106,
		5.6112,
		6.012,
		6.4128,
		6.8136,
		7.2144,
		7.6152,
		8.016,
		8.517,
		9.018,
		9.519,
	}
	resolveBaseBonus = []float64{
		0.0389,
		0.0418,
		0.0447,
		0.0486,
		0.0515,
		0.0544,
		0.0583,
		0.0622,
		0.0661,
		0.07,
		0.0739,
		0.0778,
		0.0826,
		0.0875,
		0.0923,
	}
	resolveBonus = []float64{
		0.0073,
		0.0078,
		0.0084,
		0.0091,
		0.0096,
		0.0102,
		0.0109,
		0.0116,
		0.0123,
		0.0131,
		0.0138,
		0.0145,
		0.0154,
		0.0163,
		0.0172,
	}
	resolveBonusOld = []float64{
		0.0048,
		0.0052,
		0.0056,
		0.0061,
		0.0064,
		0.0068,
		0.0073,
		0.0078,
		0.0082,
		0.0087,
		0.0092,
		0.0097,
		0.0103,
		0.0109,
		0.0115,
	}
	burstRestore = []float64{
		1.6,
		1.7,
		1.8,
		1.9,
		2,
		2.1,
		2.2,
		2.3,
		2.4,
		2.5,
		2.5,
		2.5,
		2.5,
		2.5,
		2.5,
	}
	resolveStackGain = []float64{
		0.15,
		0.16,
		0.16,
		0.17,
		0.17,
		0.18,
		0.18,
		0.19,
		0.19,
		0.2,
		0.2,
		0.2,
		0.2,
		0.2,
		0.2,
	}
)
